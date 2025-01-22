import request from '@/utils/request'

export function getMiddlewareList(params) {
  return request({
    url: '/api/v1/middleware/list',
    method: 'get',
    params
  })
}

export function createMiddleware(data) {
  return request({
    url: '/api/v1/middleware/create',
    method: 'post',
    data
  })
}

export function updateMiddleware(id, data) {
  return request({
    url: `/api/v1/middleware/${id}`,
    method: 'put',
    data
  })
}

export function deleteMiddleware(id) {
  return request({
    url: `/api/v1/middleware/${id}`,
    method: 'delete'
  })
} 