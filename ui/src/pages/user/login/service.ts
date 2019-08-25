import request from '@/utils/request';
import { FromDataType } from './index';

export async function fakeAccountLogin(params: FromDataType) {
  return request('/api/login', {
    method: 'POST',
    data: params,
  });
}

