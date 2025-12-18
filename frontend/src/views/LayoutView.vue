<template>
  <div class="layout">
    <aside class="sidebar">
      <div class="logo">人事调动系统</div>
      <nav class="menu">
        <button class="menu-item" :class="{ active: isActive('/employees') }" @click="go('/employees')">
          员工管理
        </button>
        <button class="menu-item" :class="{ active: isActive('/departments') }" @click="go('/departments')">
          部门管理
        </button>
        <button class="menu-item" :class="{ active: isActive('/transfers') }" @click="go('/transfers')">
          调动管理
        </button>
        <button class="menu-item" :class="{ active: isActive('/backup') }" @click="go('/backup')">
          系统维护
        </button>
      </nav>
      <div class="sidebar-footer">
        <div class="user-info">
          <div class="user-name">{{ userName }}</div>
        </div>
        <button class="logout-button" @click="logout">退出登录</button>
      </div>
    </aside>
    <main class="content">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { computed } from "vue"
import { useRoute, useRouter } from "vue-router"

const route = useRoute()
const router = useRouter()

const userName = computed(() => {
  const raw = localStorage.getItem("user")
  if (!raw) {
    return "未登录"
  }
  try {
    const u = JSON.parse(raw)
    return u.real_name || u.username || "管理员"
  } catch {
    return "管理员"
  }
})

const isActive = path => {
  return route.path.startsWith(path)
}

const go = path => {
  router.push(path)
}

const logout = () => {
  localStorage.removeItem("token")
  localStorage.removeItem("user")
  router.push("/login")
}
</script>
