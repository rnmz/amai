import { defineStore } from "pinia";
import { ref } from "vue";
import {
  getPost,
  getPosts,
  createPost,
  editPost,
  deletePost,
  type Post,
  type CreatePostPayload,
  type EditPostPayload,
} from "@/api/post";
import { getApiErrorMessage } from "@/api/client";

export const usePostsStore = defineStore("posts", () => {
  const posts = ref<Post[]>([]);
  const pages = ref(0);
  const currentPage = ref(1);
  const isLoading = ref(false);
  const error = ref<string | null>(null);

  async function fetchPosts(page = currentPage.value): Promise<void> {
    isLoading.value = true;
    error.value = null;
    try {
      const data = await getPosts(page);
      posts.value = data.posts;
      pages.value = data.pages;
      currentPage.value = page;
    } catch (e) {
      error.value = getApiErrorMessage(e, "Unable to load articles");
    } finally {
      isLoading.value = false;
    }
  }

  async function create(payload: CreatePostPayload): Promise<boolean> {
    isLoading.value = true;
    error.value = null;
    try {
      await createPost(payload);
      return true;
    } catch (e) {
      error.value = getApiErrorMessage(e, "Failed to create article");
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  async function update(payload: EditPostPayload): Promise<boolean> {
    isLoading.value = true;
    error.value = null;
    try {
      await editPost(payload);
      const idx = posts.value.findIndex((p) => p.id === payload.id);
      const existing = idx !== -1 ? posts.value[idx] : undefined;
      if (existing) {
        posts.value[idx] = { ...existing, ...payload };
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
      await deletePost(id);
      posts.value = posts.value.filter((p) => p.id !== id);
      return true;
    } catch (e) {
      error.value = getApiErrorMessage(e, "Failed to delete article");
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  return {
    posts,
    pages,
    currentPage,
    isLoading,
    error,
    fetchPosts,
    create,
    update,
    remove,
  };
});

export const usePostStore = defineStore("post", () => {
  const post = ref<Post | null>(null);
  const isLoading = ref(false);
  const error = ref<string | null>(null);

  async function fetchPost(id: string): Promise<void> {
    isLoading.value = true;
    error.value = null;
    try {
      post.value = await getPost(id);
    } catch (e) {
      error.value = getApiErrorMessage(e, "Article not found");
      post.value = null;
    } finally {
      isLoading.value = false;
    }
  }

  return { post, isLoading, error, fetchPost };
});