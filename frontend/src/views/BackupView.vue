<template>
  <div>
    <h2 class="page-title">系统维护与备份</h2>
    <p>点击下方按钮可以导出当前所有员工的基本信息为 CSV 文件。</p>
    <button class="primary-button" @click="exportEmployees" :disabled="downloading">
      {{ downloading ? "导出中..." : "导出员工信息" }}
    </button>
  </div>
</template>

<script setup>
import { ref } from "vue"

const downloading = ref(false)

const exportEmployees = async () => {
  downloading.value = true
  try {
    const token = localStorage.getItem("token")
    const headers = {}
    if (token) {
      headers.Authorization = "Bearer " + token
    }
    const res = await fetch("/api/backup/export", {
      headers
    })
    const blob = await res.blob()
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement("a")
    a.href = url
    a.download = "employees_backup.csv"
    document.body.appendChild(a)
    a.click()
    a.remove()
    window.URL.revokeObjectURL(url)
  } catch (e) {
    alert("导出失败")
  } finally {
    downloading.value = false
  }
}
</script>

