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
const MonitorList = () => import('@/pages/monitor/MonitorList.vue')
const MonitorEdit = () => import('@/pages/monitor/MonitorEdit.vue')
const ChannelList = () => import('@/pages/monitor/ChannelList.vue')
const ChannelEdit = () => import('@/pages/monitor/ChannelEdit.vue')

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', component: Login, meta: { public: true, locale: 'login.title' } },
    {
      path: '/',
      component: DefaultLayout,
      children: [
        { path: '', redirect: '/dashboard' },
        { path: 'dashboard', component: Dashboard, meta: { locale: 'menu.dashboard', localeSubtitle: 'dashboard.subtitle' } },
        { path: 'logs', component: Logs, meta: { locale: 'menu.logs', localeSubtitle: 'logs.subtitle' } },
        { path: 'models', component: Models, meta: { locale: 'menu.models', localeSubtitle: 'models.subtitle' } },
        { path: 'datasources', component: DataSources, meta: { locale: 'menu.datasources', localeSubtitle: 'datasource.subtitle' } },
        { path: 'datasources/new', component: NewDataSource, meta: { locale: 'datasource.add', localeSubtitle: 'datasource.newSubtitle' } },
        { path: 'datasources/new/loki', component: LokiConfig, meta: { locale: 'datasource.lokiConfig', localeSubtitle: 'datasource.configSubtitle' } },
        { path: 'datasources/new/elasticsearch', component: ElasticsearchConfig, meta: { locale: 'datasource.esConfig', localeSubtitle: 'datasource.configSubtitle' } },
        { path: 'datasources/new/victorialogs', component: VictoriaLogsConfig, meta: { locale: 'datasource.vlConfig', localeSubtitle: 'datasource.configSubtitle' } },
        { path: 'profile', component: Profile, meta: { locale: 'profile.accountInfo', localeSubtitle: 'profile.subtitle' } },
        { path: 'monitors', component: MonitorList, meta: { locale: 'menu.monitors', localeSubtitle: 'monitor.subtitle' } },
        { path: 'monitors/new', component: MonitorEdit, meta: { locale: 'monitor.newTask', localeSubtitle: 'monitor.newSubtitle' } },
        { path: 'monitors/:id', component: MonitorEdit, meta: { locale: 'monitor.editTask', localeSubtitle: 'monitor.editSubtitle' } },
        { path: 'channels', component: ChannelList, meta: { locale: 'menu.channels', localeSubtitle: 'channel.subtitle' } },
        { path: 'channels/new', component: ChannelEdit, meta: { locale: 'channel.newChannel', localeSubtitle: 'channel.newSubtitle' } },
        { path: 'channels/:id', component: ChannelEdit, meta: { locale: 'channel.editChannel', localeSubtitle: 'channel.editSubtitle' } },
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

