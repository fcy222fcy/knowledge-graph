import request from '@/utils/request'

export interface GraphNode {
  id: number
  name: string
  description: string
  document_id: number
  category?: string
  // D3 force simulation properties
  x?: number
  y?: number
  vx?: number
  vy?: number
  fx?: number | null
  fy?: number | null
}

export interface GraphEdge {
  id: number
  source: number
  target: number
  relation_type: string
  description?: string
}

export interface GraphData {
  nodes: GraphNode[]
  edges: GraphEdge[]
  summary: {
    node_count: number
    edge_count: number
  }
}

export interface GraphBuildResult {
  build_id: number
  created_points: number
  created_relations: number
  chunk_count: number
  vector_count: number
  status: string
  message: string
}

// 获取图谱数据
export function getGraphData(params?: { document_id?: number; keyword?: string; relation_type?: string }) {
  return request.get<GraphData>('/graph', { params })
}

// 从文档构建图谱
export function buildGraph(document_ids: number[]) {
  return request.post<GraphBuildResult>('/graph/build', { document_ids })
}

// 获取最近构建结果
export function getLatestBuild() {
  return request.get('/graph/build/latest')
}

// 构建历史记录
export function getBuildHistory(params?: { page?: number; size?: number }) {
  return request.get('/graph/build/history', { params })
}
