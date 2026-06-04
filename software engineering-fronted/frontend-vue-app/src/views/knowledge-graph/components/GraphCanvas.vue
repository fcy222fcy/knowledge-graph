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
  zoom = d3.zoom<SVGSVGElement, unknown>()
    .scaleExtent([0.3, 4])
    .on('zoom', (event) => g.attr('transform', event.transform))
  svg.call(zoom)

  const nodes = props.data.nodes.map(d => ({ ...d }))
  const edges = props.data.edges.map(d => ({ ...d }))

  const categories = [...new Set(nodes.map(n => n.category || '未知'))]
  const colorScale = d3.scaleOrdinal<string>().domain(categories).range(COLORS)

  simulation = d3.forceSimulation<GraphNode>(nodes)
    .force('link', d3.forceLink<GraphNode, GraphEdge>(edges).id(d => d.id).distance(120))
    .force('charge', d3.forceManyBody().strength(-300))
    .force('center', d3.forceCenter(width / 2, height / 2))
    .force('collide', d3.forceCollide().radius(30))

  const link = g.append('g')
    .selectAll('line')
    .data(edges)
    .join('line')
    .attr('stroke', '#cbd5e1')
    .attr('stroke-width', 1.5)
    .attr('stroke-opacity', 0.6)

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

  node.append('circle')
    .attr('r', 16)
    .attr('fill', d => colorScale(d.category || '未知'))
    .attr('stroke', '#fff')
    .attr('stroke-width', 2)
    .style('cursor', 'pointer')
    .on('click', (_, d) => emit('nodeClick', d))
    .on('mouseenter', function(_, d) {
      d3.select(this).transition().attr('r', 20)
      link.attr('stroke', (l: any) => {
        const source = typeof l.source === 'object' ? l.source.id : l.source
        const target = typeof l.target === 'object' ? l.target.id : l.target
        return (source === d.id || target === d.id) ? '#3b82f6' : '#e2e8f0'
      }).attr('stroke-width', (l: any) => {
        const source = typeof l.source === 'object' ? l.source.id : l.source
        const target = typeof l.target === 'object' ? l.target.id : l.target
        return (source === d.id || target === d.id) ? 3 : 1
      })
    })
    .on('mouseleave', function() {
      d3.select(this).transition().attr('r', 16)
      link.attr('stroke', '#cbd5e1').attr('stroke-width', 1.5)
    })

  node.append('text')
    .text(d => d.name)
    .attr('dy', 30)
    .attr('text-anchor', 'middle')
    .attr('font-size', 12)
    .attr('fill', '#475569')

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
    .attr('r', d => d.name.includes(keyword) ? 22 : 14)
    .attr('stroke', d => d.name.includes(keyword) ? '#f59e0b' : '#fff')
    .attr('stroke-width', d => d.name.includes(keyword) ? 3 : 2)
}

watch(() => props.data, async () => {
  await nextTick()
  renderGraph()
}, { deep: true })

watch(() => props.highlight, (val) => {
  highlightNodes(val)
})

onMounted(() => {
  renderGraph()
  if (containerRef.value) {
    const observer = new ResizeObserver(() => renderGraph())
    observer.observe(containerRef.value)
  }
})

onUnmounted(() => {
  simulation?.stop()
  simulation = null
})
</script>

<style scoped>
.graph-canvas {
  width: 100%;
  height: 500px;
  background: var(--bg-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
}
.graph-canvas svg {
  width: 100%;
  height: 100%;
}
</style>
