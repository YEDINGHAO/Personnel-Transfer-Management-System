import { createRouter, createWebHistory } from "vue-router"
import LoginView from "../views/LoginView.vue"
import RegisterView from "../views/RegisterView.vue"
import LayoutView from "../views/LayoutView.vue"
import EmployeeListView from "../views/EmployeeListView.vue"
import DepartmentListView from "../views/DepartmentListView.vue"
import TransferListView from "../views/TransferListView.vue"
import BackupView from "../views/BackupView.vue"

const routes = [
  {
    path: "/login",
    name: "login",
    component: LoginView
  },
  {
    path: "/register",
    name: "register",
    component: RegisterView
  },
  {
    path: "/",
    component: LayoutView,
    children: [
      {
        path: "",
        redirect: "/employees"
      },
      {
        path: "employees",
        name: "employees",
        component: EmployeeListView
      },
      {
        path: "departments",
        name: "departments",
        component: DepartmentListView
      },
      {
        path: "transfers",
        name: "transfers",
        component: TransferListView
      },
      {
        path: "backup",
        name: "backup",
        component: BackupView
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  if (to.path === "/login" || to.path === "/register") {
    next()
    return
  }
  const token = localStorage.getItem("token")
  if (!token) {
    next("/login")
    return
  }
  next()
})

export default router
