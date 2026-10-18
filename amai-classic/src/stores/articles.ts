import { defineStore } from "pinia";
import { ref } from "vue";
import {
  getArticle,
  getArticles,
  createArticle,
  editArticle,
  deleteArticle,
  type Article,
  type CreateArticlePayload,
  type EditArticlePayload,
} from "@/api/articles.ts";
import { getApiErrorMessage } from "@/api/client";

export const useArticlesStore = defineStore("articles", () => {
  const articles = ref<Article[]>([]);
  const pages = ref(0);
  const currentPage = ref(1);
  const isLoading = ref(false);
  const error = ref<string | null>(null);

  async function fetchArticles(page = currentPage.value): Promise<void> {
    isLoading.value = true;
    error.value = null;
    try {
      const data = await getArticles(page);
      articles.value = data.articles;
      pages.value = data.pages;
      currentPage.value = page;
    } catch (e) {
      error.value = getApiErrorMessage(e, "Unable to load articles");
    } finally {
      isLoading.value = false;
    }
  }

  async function create(payload: CreateArticlePayload): Promise<boolean> {
    isLoading.value = true;
    error.value = null;
    try {
      await createArticle(payload);
      return true;
    } catch (e) {
      error.value = getApiErrorMessage(e, "Failed to create article");
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  async function update(payload: EditArticlePayload): Promise<boolean> {
    isLoading.value = true;
    error.value = null;
    try {
      await editArticle(payload);
      const idx = articles.value.findIndex((a) => a.id === payload.id);
      const existing = idx !== -1 ? articles.value[idx] : undefined;
      if (existing) {
        articles.value[idx] = { ...existing, ...payload };
      }
      return true;
    } catch (e) {
      error.value = getApiErrorMessage(e, "Unable to edit article");
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  async function remove(id: string): Promise<boolean> {
    isLoading.value = true;
    error.value = null;
    try {
      await deleteArticle(id);
      articles.value = articles.value.filter((a) => a.id !== id);
      return true;
    } catch (e) {
      error.value = getApiErrorMessage(e, "Failed to delete article");
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  return {
    articles,
    pages,
    currentPage,
    isLoading,
    error,
    fetchArticles,
    create,
    update,
    remove,
  };
});

export const useArticleStore = defineStore("article", () => {
  const article = ref<Article | null>(null);
  const isLoading = ref(false);
  const error = ref<string | null>(null);

  async function fetchArticle(id: string): Promise<void> {
    isLoading.value = true;
    error.value = null;
    try {
      article.value = await getArticle(id);
    } catch (e) {
      error.value = getApiErrorMessage(e, "Articles not found");
      article.value = null;
    } finally {
      isLoading.value = false;
    }
  }

  return { article, isLoading, error, fetchArticle };
});
