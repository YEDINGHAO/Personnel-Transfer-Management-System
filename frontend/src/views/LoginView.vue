<template>
  <div class="auth-page">
    <div class="auth-card">
      <h1 class="title">员工人事调动管理系统</h1>
      <h2 class="subtitle">管理员登录</h2>
      <form class="form" @submit.prevent="onSubmit">
        <div class="form-item">
          <label>用户名</label>
          <input v-model="form.username" type="text" required />
        </div>
        <div class="form-item">
          <label>密码</label>
          <input v-model="form.password" type="password" required />
        </div>
        <button class="primary-button" type="submit" :disabled="loading">
          {{ loading ? "登录中..." : "登录" }}
        </button>
        <p v-if="error" class="error-text">{{ error }}</p>
        <button class="link-button" type="button" @click="goRegister">
          还没有账号？去注册
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from "vue"
import { useRouter } from "vue-router"

const router = useRouter()
const form = reactive({
  username: "",
  password: ""
})
const loading = ref(false)
const error = ref("")

const onSubmit = async () => {
  error.value = ""
  loading.value = true
  try {
    const res = await fetch("/api/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify(form)
    })
    const data = await res.json()
    if (data.code !== 0) {
      error.value = data.message || "登录失败"
      return
    }
    localStorage.setItem("token", data.data.token)
    localStorage.setItem("user", JSON.stringify(data.data.user))
    router.push("/")
  } catch (e) {
    error.value = "网络错误"
  } finally {
    loading.value = false
  }
}

const goRegister = () => {
  router.push("/register")
}
</script>
