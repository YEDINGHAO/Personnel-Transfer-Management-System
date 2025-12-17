<template>
  <div>
    <h2 class="page-title">调动管理</h2>
    <div class="toolbar">
      <div class="filters">
        <input v-model="filters.employee_id" placeholder="员工ID" />
        <select v-model="filters.status">
          <option value="">全部状态</option>
          <option value="1">待审批</option>
          <option value="2">已批准</option>
          <option value="3">已驳回</option>
        </select>
        <button class="primary-button" @click="loadTransfers">查询</button>
      </div>
      <button class="primary-button" @click="openCreate">新建调动/离退休</button>
    </div>
    <table class="table">
      <thead>
        <tr>
          <th>ID</th>
          <th>员工</th>
          <th>类型</th>
          <th>日期</th>
          <th>调出部门</th>
          <th>调入部门</th>
          <th>原因</th>
          <th>状态</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="t in transfers" :key="t.id">
          <td>{{ t.id }}</td>
          <td>{{ t.employee?.name || t.employee_id }}</td>
          <td>{{ typeText(t.type) }}</td>
          <td>{{ t.transfer_date }}</td>
          <td>{{ t.from_dept?.name || t.from_dept_id }}</td>
          <td>{{ t.to_dept?.name || t.to_dept_id }}</td>
          <td>{{ t.reason }}</td>
          <td>{{ statusText(t.status) }}</td>
          <td>
            <button
              v-if="t.status === 1"
              class="link-button"
              @click="approve(t, 2)"
            >
              通过
            </button>
            <button
              v-if="t.status === 1"
              class="link-button danger"
              @click="approve(t, 3)"
            >
              驳回
            </button>
          </td>
        </tr>
        <tr v-if="transfers.length === 0">
          <td colspan="9" class="empty-cell">暂无数据</td>
        </tr>
      </tbody>
    </table>
    <div v-if="showDialog" class="dialog-mask">
      <div class="dialog">
        <h3>新建调动或离退休申请</h3>
        <form class="form" @submit.prevent="submitForm">
          <div class="form-item">
            <label>员工ID</label>
            <input v-model.number="form.employee_id" type="number" required />
          </div>
          <div class="form-item">
            <label>调动类型</label>
            <select v-model.number="form.type" required>
              <option :value="1">部门调动</option>
              <option :value="2">职位调动</option>
              <option :value="3">离退休</option>
            </select>
          </div>
          <div class="form-item">
            <label>调动日期</label>
            <input v-model="form.transfer_date" type="date" required />
          </div>
          <div class="form-item">
            <label>调出部门ID</label>
            <input v-model.number="form.from_dept_id" type="number" />
          </div>
          <div class="form-item">
            <label>调入部门ID</label>
            <input v-model.number="form.to_dept_id" type="number" />
          </div>
          <div class="form-item">
            <label>调动原因</label>
            <textarea v-model="form.reason" rows="2" />
          </div>
          <div class="dialog-actions">
            <button type="button" @click="closeDialog">取消</button>
            <button class="primary-button" type="submit">提交申请</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from "vue"

const transfers = ref([])
const filters = reactive({
  employee_id: "",
  status: ""
})

const showDialog = ref(false)
const form = reactive({
  employee_id: null,
  type: 1,
  transfer_date: "",
  from_dept_id: null,
  to_dept_id: null,
  reason: ""
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

const typeText = v => {
  if (v === 1) {
    return "部门调动"
  }
  if (v === 2) {
    return "职位调动"
  }
  if (v === 3) {
    return "离退休"
  }
  return String(v)
}

const statusText = v => {
  if (v === 1) {
    return "待审批"
  }
  if (v === 2) {
    return "已批准"
  }
  if (v === 3) {
    return "已驳回"
  }
  if (v === 4) {
    return "已完成"
  }
  return String(v)
}

const loadTransfers = async () => {
  const params = new URLSearchParams()
  if (filters.employee_id) {
    params.append("employee_id", filters.employee_id)
  }
  if (filters.status) {
    params.append("status", filters.status)
  }
  const res = await fetch("/api/transfers?" + params.toString(), {
    headers: authHeaders()
  })
  const data = await res.json()
  if (data.code === 0) {
    transfers.value = data.data
  }
}

const openCreate = () => {
  Object.assign(form, {
    employee_id: null,
    type: 1,
    transfer_date: "",
    from_dept_id: null,
    to_dept_id: null,
    reason: ""
  })
  showDialog.value = true
}

const closeDialog = () => {
  showDialog.value = false
}

const submitForm = async () => {
  const payload = {
    employee_id: form.employee_id,
    type: form.type,
    transfer_date: form.transfer_date,
    from_dept_id: form.from_dept_id || 0,
    to_dept_id: form.to_dept_id || 0,
    reason: form.reason
  }
  const res = await fetch("/api/transfers", {
    method: "POST",
    headers: authHeaders(),
    body: JSON.stringify(payload)
  })
  const data = await res.json()
  if (data.code === 0) {
    showDialog.value = false
    loadTransfers()
  } else {
    alert(data.message || "提交失败")
  }
}

const approve = async (t, status) => {
  const ok = confirm(status === 2 ? "确认通过该调动申请吗" : "确认驳回该调动申请吗")
  if (!ok) {
    return
  }
  const userRaw = localStorage.getItem("user")
  let approverId = 0
  if (userRaw) {
    try {
      const u = JSON.parse(userRaw)
      approverId = u.id
    } catch {
      approverId = 0
    }
  }
  const res = await fetch(`/api/transfers/${t.id}/approve`, {
    method: "PUT",
    headers: authHeaders(),
    body: JSON.stringify({
      status,
      approver_id: approverId
    })
  })
  const data = await res.json()
  if (data.code === 0) {
    loadTransfers()
  } else {
    alert(data.message || "操作失败")
  }
}

onMounted(() => {
  loadTransfers()
})
</script>

