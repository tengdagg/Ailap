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
const VictoriaLogsConfig = () => import('@/pages/datasource/VictoriaLogsConfig.vue')
const Profile = () => import('@/pages/Profile.vue')

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', component: Login, meta: { public: true, title: '登录' } },
    {
      path: '/',
      component: DefaultLayout,
      children: [
        { path: '', redirect: '/dashboard' },
        { path: 'dashboard', component: Dashboard, meta: { title: '仪表盘', pageSubtitle: '系统概览和监控面板' } },
        { path: 'logs', component: Logs, meta: { title: '日志', pageSubtitle: '日志查询和分析工具' } },
        { path: 'models', component: Models, meta: { title: '模型', pageSubtitle: 'AI 模型配置和管理' } },
        { path: 'datasources', component: DataSources, meta: { title: '数据源', pageSubtitle: '维护日志/数据系统连接' } },
        { path: 'datasources/new', component: NewDataSource, meta: { title: '添加数据源', pageSubtitle: '选择数据源类型' } },
        { path: 'datasources/new/loki', component: LokiConfig, meta: { title: 'Loki 配置', pageSubtitle: '配置 Loki 数据源连接信息' } },
        { path: 'datasources/new/elasticsearch', component: ElasticsearchConfig, meta: { title: 'Elasticsearch 配置', pageSubtitle: '配置 Elasticsearch 数据源连接信息' } },
        { path: 'datasources/new/victorialogs', component: VictoriaLogsConfig, meta: { title: 'VictoriaLogs 配置', pageSubtitle: '配置 VictoriaLogs 数据源连接信息' } },
        { path: 'profile', component: Profile, meta: { title: '个人中心', pageSubtitle: '账户资料与安全设置' } },
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

