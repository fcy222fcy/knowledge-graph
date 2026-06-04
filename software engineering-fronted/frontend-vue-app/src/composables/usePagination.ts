import { reactive } from 'vue'

export function usePagination(defaultSize = 10) {
  const pagination = reactive({
    page: 1,
    size: defaultSize,
    total: 0
  })

  const handlePageChange = (page: number) => {
    pagination.page = page
  }

  const handleSizeChange = (size: number) => {
    pagination.size = size
    pagination.page = 1
  }

  const resetPage = () => {
    pagination.page = 1
  }

  return {
    pagination,
    handlePageChange,
    handleSizeChange,
    resetPage
  }
}
