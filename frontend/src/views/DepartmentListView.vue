<template>
  <div>
    <h2 class="page-title">部门管理</h2>
    <div class="toolbar">
      <button class="primary-button" @click="openCreate">新增部门</button>
    </div>
    <table class="table">
      <thead>
        <tr>
          <th>ID</th>
          <th>部门编号</th>
          <th>部门名称</th>
          <th>主管员工ID</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="d in departments" :key="d.id">
          <td>{{ d.id }}</td>
          <td>{{ d.dept_no }}</td>
          <td>{{ d.name }}</td>
          <td>{{ d.manager_id }}</td>
          <td>
            <button class="link-button" @click="edit(d)">编辑</button>
            <button class="link-button danger" @click="remove(d)">删除</button>
          </td>
        </tr>
        <tr v-if="departments.length === 0">
          <td colspan="5" class="empty-cell">暂无数据</td>
        </tr>
      </tbody>
    </table>
    <div v-if="showDialog" class="dialog-mask">
      <div class="dialog">
        <h3>{{ editing ? "编辑部门" : "新增部门" }}</h3>
        <form class="form" @submit.prevent="submitForm">
          <div class="form-item">
            <label>部门编号</label>
            <input v-model="form.dept_no" :disabled="editing" required />
          </div>
          <div class="form-item">
            <label>部门名称</label>
            <input v-model="form.name" required />
          </div>
          <div class="form-item">
            <label>主管员工ID</label>
            <input v-model.number="form.manager_id" type="number" min="0" />
          </div>
          <div class="dialog-actions">
            <button type="button" @click="closeDialog">取消</button>
            <button class="primary-button" type="submit">保存</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from "vue"

const departments = ref([])
const showDialog = ref(false)
const editing = ref(false)
const currentId = ref(null)
const form = reactive({
  dept_no: "",
  name: "",
  manager_id: 0
})

const authHeaders = () => {
  const token = localStorage.getItem("token")
  const headers = {}
  if (token) {
    headers.Authorization = "Bearer " + token
  }
  headers["Content-Type"] = "application/json"
  return headers
}

const loadDepartments = async () => {
  const res = await fetch("/api/departments", {
    headers: authHeaders()
  })
  const data = await res.json()
  if (data.code === 0) {
    departments.value = data.data || []
  }
}

const openCreate = () => {
  editing.value = false
  currentId.value = null
  Object.assign(form, {
    dept_no: "",
    name: "",
    manager_id: 0
  })
  showDialog.value = true
}

const edit = d => {
  editing.value = true
  currentId.value = d.id
  Object.assign(form, {
    dept_no: d.dept_no,
    name: d.name,
    manager_id: d.manager_id || 0
  })
  showDialog.value = true
}

const closeDialog = () => {
  showDialog.value = false
}

const submitForm = async () => {
  const payload = {
    dept_no: form.dept_no,
    name: form.name,
    manager_id: form.manager_id || 0
  }
  let url = "/api/departments"
  let method = "POST"
  if (editing.value && currentId.value) {
    url = "/api/departments/" + currentId.value
    method = "PUT"
  }
  const res = await fetch(url, {
    method,
    headers: authHeaders(),
    body: JSON.stringify(payload)
  })
  const data = await res.json()
  if (data.code === 0) {
    showDialog.value = false
    loadDepartments()
  } else {
    alert(data.message || "保存失败")
  }
}

const remove = async d => {
  if (!confirm("确认删除该部门吗")) {
    return
  }
  const res = await fetch("/api/departments/" + d.id, {
    method: "DELETE",
    headers: authHeaders()
  })
  const data = await res.json()
  if (data.code === 0) {
    loadDepartments()
  } else {
    alert(data.message || "删除失败：该部门可能存在关联的调动记录")
  }
}

onMounted(() => {
  loadDepartments()
})
</script>
