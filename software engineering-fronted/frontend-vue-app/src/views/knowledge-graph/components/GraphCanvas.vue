<template>
  <div ref="containerRef" class="graph-canvas" v-loading="loading">
    <svg ref="svgRef"></svg>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted, nextTick } from 'vue'
import * as d3 from 'd3'
import type { GraphData, GraphNode, GraphEdge } from '@/types/graph'

const props = defineProps<{
  data: GraphData | null
  loading: boolean
  highlight: string
}>()

const emit = defineEmits<{
  nodeClick: [node: GraphNode]
}>()

const containerRef = ref<HTMLDivElement>()
const svgRef = ref<SVGSVGElement>()
let simulation: d3.Simulation<GraphNode, GraphEdge> | null = null
let zoom: d3.ZoomBehavior<SVGSVGElement, unknown> | null = null

const COLORS = ['#3b82f6', '#10b981', '#f59e0b', '#f43f5e', '#8b5cf6', '#06b6d4', '#ec4899']

const renderGraph = () => {
  if (!props.data || !svgRef.value || !containerRef.value) return

  d3.select(svgRef.value).selectAll('*').remove()

  const width = containerRef.value.clientWidth
  const height = containerRef.value.clientHeight

  const svg = d3.select(svgRef.value)
    .attr('width', width)
    .attr('height', height)

  const g = svg.append('g')
  let currentScale = 1

  const nodes = (props.data.nodes || []).map(d => ({ ...d }))
  const edges = (props.data.edges || []).map(d => ({ ...d }))

  const categories = [...new Set(nodes.map(n => n.category || '未知'))]
  const colorScale = d3.scaleOrdinal<string>().domain(categories).range(COLORS)

  // 根据节点度数计算大小
  const degreeMap = new Map<string, number>()
  edges.forEach(e => {
    const sourceId = typeof e.source === 'object' ? (e.source as any).id : e.source
    const targetId = typeof e.target === 'object' ? (e.target as any).id : e.target
    degreeMap.set(sourceId, (degreeMap.get(sourceId) || 0) + 1)
    degreeMap.set(targetId, (degreeMap.get(targetId) || 0) + 1)
  })
  const maxDegree = Math.max(...Array.from(degreeMap.values()), 1)

  simulation = d3.forceSimulation<GraphNode>(nodes)
    .force('link', d3.forceLink<GraphNode, GraphEdge>(edges).id(d => d.id).distance(80).strength(0.3))
    .force('charge', d3.forceManyBody().strength(-100))
    .force('center', d3.forceCenter(width / 2, height / 2))
    .force('collide', d3.forceCollide().radius(25))

  const link = g.append('g')
    .selectAll('line')
    .data(edges)
    .join('line')
    .attr('stroke', '#cbd5e1')
    .attr('stroke-width', 0.5)
    .attr('stroke-opacity', 0.3)

  const node = g.append('g')
    .selectAll('g')
    .data(nodes)
    .join('g')
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    .call(d3.drag<any, GraphNode>()
      .on('start', (event, d) => {
        if (!event.active) simulation!.alphaTarget(0.3).restart()
        d.fx = d.x
        d.fy = d.y
      })
      .on('drag', (event, d) => {
        d.fx = event.x
        d.fy = event.y
      })
      .on('end', (event, d) => {
        if (!event.active) simulation!.alphaTarget(0)
        d.fx = null
        d.fy = null
      })
    )

  // 根据缩放级别调整标签显示和线条样式 - Obsidian 风格
  const updateLabels = (scale: number) => {
    const sortedDegrees = Array.from(degreeMap.values()).sort((a, b) => b - a)
    const coreThreshold = sortedDegrees[Math.floor(sortedDegrees.length * 0.2)] || 0

    // 根据缩放级别调整线条样式 - 缩小时更黑更明显
    link
      .attr('stroke', scale < 0.5 ? '#94a3b8' : scale < 0.8 ? '#a1a1aa' : '#cbd5e1')
      .attr('stroke-width', scale < 0.3 ? 1.5 : scale < 0.5 ? 1.2 : scale < 0.8 ? 0.8 : 0.5)
      .attr('stroke-opacity', scale < 0.3 ? 0.8 : scale < 0.5 ? 0.6 : scale < 0.8 ? 0.5 : 0.3)

    node.select('text')
      .attr('opacity', (d: any) => {
        const degree = degreeMap.get(d.id) || 0
        // 核心节点在中等缩放时就显示
        if (degree >= coreThreshold && degree > 3) {
          return scale > 0.6 ? 0.9 : 0
        }
        // 普通节点放大后才显示
        return scale > 1.2 ? 0.7 : 0
      })
      .attr('font-size', (d: any) => {
        const degree = degreeMap.get(d.id) || 0
        if (degree >= coreThreshold && degree > 3) return '10px'
        return '8px'
      })
      .attr('font-weight', (d: any) => {
        const degree = degreeMap.get(d.id) || 0
        return degree >= coreThreshold && degree > 3 ? 600 : 400
      })
  }

  zoom = d3.zoom<SVGSVGElement, unknown>()
    .scaleExtent([0.1, 6])
    .on('zoom', (event) => {
      g.attr('transform', event.transform)
      currentScale = event.transform.k
      updateLabels(currentScale)
    })
  svg.call(zoom)

  // 初始居中
  const initialTransform = d3.zoomIdentity.translate(width / 2, height / 2).scale(0.7).translate(-width / 2, -height / 2)
  svg.call(zoom!.transform, initialTransform)

  // 节点圆圈 - 统一小尺寸
  node.append('circle')
    .attr('r', 5)
    .attr('fill', d => colorScale(d.category || '未知'))
    .attr('stroke', '#fff')
    .attr('stroke-width', 1)
    .style('cursor', 'pointer')
    .on('click', (_, d) => emit('nodeClick', d))
    .on('mouseenter', function(_, d) {
      d3.select(this).transition().duration(200).attr('r', 7)
      // 高亮关联边
      link.transition().duration(200)
        .attr('stroke', (l: any) => {
          const source = typeof l.source === 'object' ? l.source.id : l.source
          const target = typeof l.target === 'object' ? l.target.id : l.target
          return (source === d.id || target === d.id) ? colorScale(d.category || '未知') : '#f1f5f9'
        })
        .attr('stroke-width', (l: any) => {
          const source = typeof l.source === 'object' ? l.source.id : l.source
          const target = typeof l.target === 'object' ? l.target.id : l.target
          return (source === d.id || target === d.id) ? 2.5 : 0.5
        })
        .attr('stroke-opacity', (l: any) => {
          const source = typeof l.source === 'object' ? l.source.id : l.source
          const target = typeof l.target === 'object' ? l.target.id : l.target
          return (source === d.id || target === d.id) ? 1 : 0.2
        })
      // 显示关联节点标签
      node.select('text').transition().duration(200)
        .attr('opacity', (n: any) => {
          if (n.id === d.id) return 1
          const isConnected = edges.some((e: any) => {
            const s = typeof e.source === 'object' ? e.source.id : e.source
            const t = typeof e.target === 'object' ? e.target.id : e.target
            return (s === d.id && t === n.id) || (t === d.id && s === n.id)
          })
          return isConnected ? 1 : 0.15
        })
    })
    .on('mouseleave', function() {
      d3.select(this).transition().duration(200).attr('r', 5)
      link.transition().duration(200)
        .attr('stroke', '#e2e8f0')
        .attr('stroke-width', 1)
        .attr('stroke-opacity', 0.5)
      updateLabels(currentScale)
    })

  // 标签 - Obsidian 风格：缩小隐藏，放大显示
  node.append('text')
    .text(d => d.name)
    .attr('dy', 16)
    .attr('text-anchor', 'middle')
    .attr('font-size', '8px')
    .attr('font-weight', 400)
    .attr('fill', '#475569')
    .attr('pointer-events', 'none')
    .attr('opacity', 0)

  simulation.on('tick', () => {
    link
      .attr('x1', d => (d.source as any).x)
      .attr('y1', d => (d.source as any).y)
      .attr('x2', d => (d.target as any).x)
      .attr('y2', d => (d.target as any).y)
    node.attr('transform', d => `translate(${d.x},${d.y})`)
  })
}

const highlightNodes = (keyword: string) => {
  if (!svgRef.value || !keyword) return
  d3.select(svgRef.value).selectAll<SVGCircleElement, GraphNode>('circle')
    .transition()
    .attr('stroke', d => d.name.includes(keyword) ? '#f59e0b' : '#fff')
    .attr('stroke-width', d => d.name.includes(keyword) ? 3 : 1.5)
}

watch(() => props.data, async () => {
  await nextTick()
  renderGraph()
}, { deep: true })

watch(() => props.highlight, (val) => {
  highlightNodes(val)
})

let resizeObserver: ResizeObserver | null = null

onMounted(() => {
  renderGraph()
  if (containerRef.value) {
    resizeObserver = new ResizeObserver(() => renderGraph())
    resizeObserver.observe(containerRef.value)
  }
})

onUnmounted(() => {
  resizeObserver?.disconnect()
  resizeObserver = null
  simulation?.stop()
  simulation = null
})
</script>

<style scoped>
.graph-canvas {
  width: 100%;
  height: 100%;
  background: #fff;
  overflow: hidden;
}
.graph-canvas svg {
  width: 100%;
  height: 100%;
}
</style>
