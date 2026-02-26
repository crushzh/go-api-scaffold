import { request } from '@umijs/max';

/** List examples with pagination */
export async function getExamples(params: {
  page?: number;
  page_size?: number;
  keyword?: string;
}) {
  return request('/examples', {
    method: 'GET',
    params,
  });
}

/** Get example by ID */
export async function getExample(id: number) {
  return request(`/examples/${id}`, {
    method: 'GET',
  });
}

/** Create example */
export async function createExample(data: {
  name: string;
  description?: string;
}) {
  return request('/examples', {
    method: 'POST',
    data,
  });
}

/** Update example */
export async function updateExample(
  id: number,
  data: { name?: string; description?: string; status?: string },
) {
  return request(`/examples/${id}`, {
    method: 'PUT',
    data,
  });
}

/** Delete example */
export async function deleteExample(id: number) {
  return request(`/examples/${id}`, {
    method: 'DELETE',
  });
}
