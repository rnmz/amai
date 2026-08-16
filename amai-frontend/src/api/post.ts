import { api } from "./client";
import type { ApiMessage } from "./client";

export interface Post {
  id: string;
  title: string;
  poster_id: string;
  created_at: string;
  updated_at: string | null; 
  body: string;
}

export interface PostListResponse {
  posts: Post[];
  pages: number;
}

export interface CreatePostPayload {
  title: string;
  poster_id: string;
  body: string;
}

export interface EditPostPayload extends CreatePostPayload {
  id: string;
}

export async function getPost(id: string): Promise<Post> {
  const { data } = await api.get<Post>("/post/get", { params: { id } });
  return data;
}

export async function getPosts(page: number): Promise<PostListResponse> {
  const { data } = await api.get<PostListResponse>("/post/all", {
    params: { page },
  });
  return data;
}

export async function createPost(
  payload: CreatePostPayload
): Promise<ApiMessage> {
  const { data } = await api.post<ApiMessage>("/admin/post/create", payload);
  return data;
}

export async function editPost(payload: EditPostPayload): Promise<ApiMessage> {
  const { data } = await api.put<ApiMessage>("/admin/post/edit", payload);
  return data;
}

export async function deletePost(id: string): Promise<ApiMessage> {
  const { data } = await api.delete<ApiMessage>("/admin/post/delete", {
    params: { id },
  });
  return data;
}