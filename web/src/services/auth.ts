import { request } from '@umijs/max';

export async function login(data: { username: string; password: string }) {
  return request('/auth/login', {
    method: 'POST',
    data,
  });
}

export async function refreshToken() {
  return request('/auth/refresh', {
    method: 'POST',
  });
}
