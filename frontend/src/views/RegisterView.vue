<template>
  <div class="auth-page">
    <div class="auth-card">
      <h1 class="title">员工人事调动管理系统</h1>
      <h2 class="subtitle">管理员注册</h2>
      <form class="form" @submit.prevent="onSubmit">
        <div class="form-item">
          <label>用户名</label>
          <input v-model="form.username" type="text" required />
        </div>
        <div class="form-item">
          <label>密码</label>
          <input v-model="form.password" type="password" required />
        </div>
        <div class="form-item">
          <label>确认密码</label>
          <input v-model="form.confirmPassword" type="password" required />
        </div>
        <div class="form-item">
          <label>姓名</label>
          <input v-model="form.realName" type="text" />
        </div>
        <div class="form-item">
          <label>邮箱</label>
          <input v-model="form.email" type="email" />
        </div>
        <div class="form-item">
          <label>手机号</label>
          <input v-model="form.phone" type="text" />
        </div>
        <button class="primary-button" type="submit" :disabled="loading">
          {{ loading ? "注册中..." : "注册" }}
        </button>
        <p v-if="error" class="error-text">{{ error }}</p>
        <button class="link-button" type="button" @click="goLogin">
          已有账号？去登录
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
  password: "",
  confirmPassword: "",
  realName: "",
  email: "",
  phone: ""
})
const loading = ref(false)
const error = ref("")

const onSubmit = async () => {
  error.value = ""
  if (!form.username || !form.password) {
    error.value = "用户名和密码为必填项"
    return
  }
  if (form.password !== form.confirmPassword) {
    error.value = "两次输入的密码不一致"
    return
  }

  loading.value = true
  try {
    const res = await fetch("/api/register", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        username: form.username,
        password: form.password,
        real_name: form.realName,
        email: form.email,
        phone: form.phone
      })
    })
    const data = await res.json()
    if (data.code !== 0) {
      error.value = data.message || "注册失败"
      return
    }
    router.push("/login")
  } catch (e) {
    error.value = "网络错误"
  } finally {
    loading.value = false
  }
}

const goLogin = () => {
  router.push("/login")
}
</script>

