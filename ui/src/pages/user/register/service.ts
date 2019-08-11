import request from '@/utils/request';
import { UserRegisterParams } from './index';

export async function sendRegister(params: UserRegisterParams) {
  return request('/api/register', {
    method: 'POST',
    data: params,
  });
}
