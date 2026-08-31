import { api } from "./client";
import type { ApiMessage } from "./client";

export interface Article {
  id: string;
  title: string;
  poster_id: string;
  created_at: string;
  updated_at: string | null;
  body: string;
}

export interface ArticlesListResponse {
  articles: Article[];
  pages: number;
}

export interface CreateArticlePayload {
  title: string;
  poster_id: string;
  body: string;
}

export interface EditArticlePayload extends CreateArticlePayload {
  id: string;
}

export async function getArticle(id: string): Promise<Article> {
  const { data } = await api.get<Article>("/article/get", { params: { id } });
  return data;
}

export async function getArticles(page: number): Promise<ArticlesListResponse> {
  const { data } = await api.get<ArticlesListResponse>("/article/all", {
    params: { page },
  });
  return data;
}

export async function createArticle(
  payload: CreateArticlePayload
): Promise<ApiMessage> {
  const { data } = await api.post<ApiMessage>("/admin/article/create", payload);
  return data;
}

export async function editArticle(payload: EditArticlePayload): Promise<ApiMessage> {
  const { data } = await api.put<ApiMessage>("/admin/article/edit", payload);
  return data;
}

export async function deleteArticle(id: string): Promise<ApiMessage> {
  const { data } = await api.delete<ApiMessage>("/admin/article/delete", {
    params: { id },
  });
  return data;
}
