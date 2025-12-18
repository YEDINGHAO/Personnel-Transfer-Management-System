<template>
  <div>
    <h2 class="page-title">员工管理</h2>
    <div class="toolbar">
      <div class="filters">
        <input v-model="filters.name" placeholder="姓名" />
        <select v-model="filters.status">
          <option value="">全部状态</option>
          <option v-for="s in statusOptions" :key="s.value" :value="s.value">
            {{ s.label }}
          </option>
        </select>
        <input v-model="filters.department" placeholder="部门" />
        <button class="primary-button" @click="loadEmployees">查询</button>
      </div>
      <button class="primary-button" @click="openCreate">新增员工</button>
    </div>
    <table class="table">
      <thead>
        <tr>
          <th>编号</th>
          <th>姓名</th>
          <th>状态</th>
          <th>部门</th>
          <th>职位</th>
          <th>职称</th>
          <th>入职日期</th>
          <th>电话</th>
          <th>邮箱</th>
          <th>地址</th>
          <th>备注</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="emp in employees" :key="emp.id">
          <td>{{ emp.employee_id }}</td>
          <td>{{ emp.name }}</td>
          <td>{{ emp.status_text }}</td>
          <td>{{ emp.department }}</td>
          <td>{{ emp.position }}</td>
          <td>{{ emp.job_title }}</td>
          <td>{{ formatDate(emp.arrival_date) }}</td>
          <td>{{ emp.phone }}</td>
          <td>{{ emp.email }}</td>
          <td>{{ emp.address }}</td>
          <td>{{ emp.remark }}</td>
          <td>
            <button class="link-button" @click="edit(emp)">编辑</button>
            <button class="link-button danger" @click="remove(emp)">删除</button>
          </td>
        </tr>
        <tr v-if="employees.length === 0">
          <td colspan="12" class="empty-cell">暂无数据</td>
        </tr>
      </tbody>
    </table>
    <div class="pagination">
      <button :disabled="page === 1" @click="changePage(page - 1)">上一页</button>
      <span>第 {{ page }} 页，共 {{ totalPages }} 页</span>
      <button :disabled="page === totalPages || totalPages === 0" @click="changePage(page + 1)">下一页</button>
    </div>
    <div v-if="showDialog" class="dialog-mask">
      <div class="dialog">
        <h3>{{ editing ? "编辑员工" : "新增员工" }}</h3>
        <form class="form" @submit.prevent="submitForm">
          <div class="form-item">
            <label>员工编号</label>
            <input v-model="form.employee_id" :disabled="editing" required />
          </div>
          <div class="form-item">
            <label>姓名</label>
            <input v-model="form.name" required />
          </div>
          <div class="form-item">
            <label>状态</label>
            <select v-model.number="form.status" required>
              <option v-for="s in statusOptions" :key="s.value" :value="s.value">
                {{ s.label }}
              </option>
            </select>
          </div>
          <div class="form-item">
            <label>入职日期</label>
            <input v-model="form.arrival_date" type="date" required />
          </div>
          <div class="form-item">
            <label>部门</label>
            <select v-model="form.department">
              <option value="">请选择部门</option>
              <option v-for="d in departments" :key="d.id" :value="d.name">
                {{ d.dept_no }} - {{ d.name }}
              </option>
            </select>
          </div>
          <div class="form-item">
            <label>职位</label>
            <input v-model="form.position" />
          </div>
          <div class="form-item">
            <label>电话</label>
            <input v-model="form.phone" />
          </div>
          <div class="form-item">
            <label>邮箱</label>
            <input v-model="form.email" type="email" />
          </div>
          <div class="form-item">
            <label>地址</label>
            <input v-model="form.address" />
          </div>
          <div class="form-item">
            <label>备注</label>
            <textarea v-model="form.remark" rows="2" />
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
import { onMounted, reactive, ref, computed } from "vue"

const employees = ref([])
const departments = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const filters = reactive({
  name: "",
  status: "",
  department: ""
})

const statusOptions = [
  { value: 1, label: "在职" },
  { value: 2, label: "兼职" },
  { value: 3, label: "试用" },
  { value: 4, label: "离职" },
  { value: 5, label: "返聘" },
  { value: 6, label: "退休" }
]

const showDialog = ref(false)
const editing = ref(false)
const currentId = ref(null)
const form = reactive({
  employee_id: "",
  name: "",
  status: 1,
  arrival_date: "",
  job_title: "",
  position: "",
  department: "",
  phone: "",
  email: "",
  address: "",
  remark: ""
})

const todayString = () => {
  const d = new Date()
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, "0")
  const day = String(d.getDate()).padStart(2, "0")
  return `${y}-${m}-${day}`
}

const totalPages = computed(() => {
  if (pageSize.value === 0) {
    return 0
  }
  return Math.ceil(total.value / pageSize.value)
})

const authHeaders = () => {
  const token = localStorage.getItem("token")
  const headers = {
    "Content-Type": "application/json"
  }
  if (token) {
    headers.Authorization = "Bearer " + token
  }
  return headers
}

const formatDate = v => {
  if (!v) {
    return ""
  }
  const s = String(v)
  return s.length >= 10 ? s.slice(0, 10) : s
}

const loadEmployees = async () => {
  const params = new URLSearchParams()
  params.append("page", String(page.value))
  params.append("page_size", String(pageSize.value))
  if (filters.name) {
    params.append("name", filters.name)
  }
  if (filters.status) {
    params.append("status", String(filters.status))
  }
  if (filters.department) {
    params.append("department", filters.department)
  }
  const res = await fetch("/api/employees?" + params.toString(), {
    headers: authHeaders()
  })
  const data = await res.json()
  if (data.code === 0) {
    employees.value = data.data.items
    total.value = data.data.total
  }
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

const changePage = p => {
  page.value = p
  loadEmployees()
}

const openCreate = () => {
  editing.value = false
  currentId.value = null
  Object.assign(form, {
    employee_id: "",
    name: "",
    status: 1,
    arrival_date: todayString(),
    job_title: "",
    position: "",
    department: "",
    phone: "",
    email: "",
    address: "",
    remark: ""
  })
  showDialog.value = true
}

const edit = emp => {
  editing.value = true
  currentId.value = emp.id
  Object.assign(form, {
    employee_id: emp.employee_id,
    name: emp.name,
    status: emp.status,
    arrival_date: emp.arrival_date,
    job_title: emp.job_title,
    position: emp.position,
    department: emp.department,
    phone: emp.phone,
    email: emp.email,
    address: emp.address,
    remark: emp.remark
  })
  showDialog.value = true
}

const closeDialog = () => {
  showDialog.value = false
}

const submitForm = async () => {
  const payload = { ...form }
  let url = "/api/employees"
  let method = "POST"
  if (editing.value && currentId.value) {
    url = "/api/employees/" + currentId.value
    method = "PUT"
    delete payload.employee_id
  }
  const res = await fetch(url, {
    method,
    headers: authHeaders(),
    body: JSON.stringify(payload)
  })
  const data = await res.json()
  if (data.code === 0) {
    showDialog.value = false
    loadEmployees()
  } else {
    alert(data.message || "保存失败")
  }
}

const remove = async emp => {
  if (!confirm("确认删除该员工吗")) {
    return
  }
  const res = await fetch("/api/employees/" + emp.id, {
    method: "DELETE",
    headers: authHeaders()
  })
  const data = await res.json()
  if (data.code === 0) {
    loadEmployees()
  } else {
    alert(data.message || "删除失败")
  }
}

onMounted(() => {
  loadEmployees()
  loadDepartments()
})
</script>
