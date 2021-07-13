import axios, { Method } from 'axios';

export const baseURL = process.env.REACT_APP_API_URL;

export const request = axios.create({
  baseURL
});

export async function authCheck(param) {
  const { data } = await request({
    ...param
  });
  return data;
}
