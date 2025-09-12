import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/store/auth'

const DefaultLayout = () => import('@/layouts/DefaultLayout.vue')
const Login = () => import('@/pages/Login.vue')
const Dashboard = () => import('@/pages/Dashboard.vue')
const Logs = () => import('@/pages/Logs.vue')
const Models = () => import('@/pages/Models.vue')
const DataSources = () => import('@/pages/DataSources.vue')
const NewDataSource = () => import('@/pages/datasource/NewDataSource.vue')
const LokiConfig = () => import('@/pages/datasource/LokiConfig.vue')
const ElasticsearchConfig = () => import('@/pages/datasource/ElasticsearchConfig.vue')

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', component: Login, meta: { public: true, title: '登录' } },
    {
      path: '/',
      component: DefaultLayout,
      children: [
        { path: '', redirect: '/dashboard' },
        { path: 'dashboard', component: Dashboard, meta: { title: '仪表盘' } },
        { path: 'logs', component: Logs, meta: { title: '日志' } },
        { path: 'models', component: Models, meta: { title: '模型' } },
        { path: 'datasources', component: DataSources, meta: { title: '数据源' } },
        { path: 'datasources/new', component: NewDataSource, meta: { title: '新建数据源' } },
        { path: 'datasources/new/loki', component: LokiConfig, meta: { title: 'Loki 配置' } },
        { path: 'datasources/new/elasticsearch', component: ElasticsearchConfig, meta: { title: 'Elasticsearch 配置' } },
      ],
    },
  ],
})

router.beforeEach((to, from, next) => {
  const auth = useAuthStore()
  if (!to.meta.public && !auth.isAuthenticated) {
    return next('/login')
  }
  next()
})

export default router

