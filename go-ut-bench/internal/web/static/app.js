function app() {
  return {
    page: 'dashboard',
    nav: [
      { id:'dashboard', icon:'<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="7" height="7" rx="1"/><rect x="14" y="3" width="7" height="7" rx="1"/><rect x="3" y="14" width="7" height="7" rx="1"/><rect x="14" y="14" width="7" height="7" rx="1"/></svg>', label:'总览' },
      { id:'new-run',   icon:'<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="5 3 19 12 5 21 5 3"/></svg>', label:'新建任务' },
      { id:'runs',      icon:'<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="8" y1="6" x2="21" y2="6"/><line x1="8" y1="12" x2="21" y2="12"/><line x1="8" y1="18" x2="21" y2="18"/><line x1="3" y1="6" x2="3.01" y2="6"/><line x1="3" y1="12" x2="3.01" y2="12"/><line x1="3" y1="18" x2="3.01" y2="18"/></svg>', label:'任务列表' },
      { id:'automations', icon:'<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="9"/><path d="M12 7v5l3 2"/><path d="M4 4l3 3"/><path d="M20 4l-3 3"/></svg>', label:'定时任务' },
      { id:'agents',    icon:'<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 3l7 4v5c0 4.5-3 7.5-7 9-4-1.5-7-4.5-7-9V7l7-4z"/><path d="M9 12h6"/><path d="M12 9v6"/></svg>', label:'Agent 接入' },
      { id:'database',  icon:'<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><ellipse cx="12" cy="5" rx="8" ry="3"/><path d="M4 5v6c0 1.7 3.6 3 8 3s8-1.3 8-3V5"/><path d="M4 11v6c0 1.7 3.6 3 8 3s8-1.3 8-3v-6"/></svg>', label:'数据库' },
      { id:'environment', icon:'<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/><path d="M9 12l2 2 4-5"/></svg>', label:'环境检查' },
      { id:'models',    icon:'<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M12 1v6m0 10v6m11-11h-6m-10 0H1m15.5-6.5l-4.25 4.25M7.75 16.25L3.5 20.5M20.5 20.5l-4.25-4.25M7.75 7.75L3.5 3.5"/></svg>', label:'模型管理' },
    ],
    config: null,
    runs: [],
    runsLoading: false,
    runFilter: '',
    statusFilter: '',
    // 评测资产（磁盘扫描视图，区别于内存中的 runs / DB 入库的 db*）
    assets: [],
    assetsLoading: false,
    assetFilter: '',
    assetQualityFilter: '', // '' | 'ok' | 'partial' | 'dirty'
    assetLanguageFilter: '',
    assetOnlyHasReport: false,
    assetOnlyIngested: false,
    selectedAssetIds: [],
    agentSubjectQuery: '',
    agentWizardOpen: false,
    agentSaving: false,
    agentFormError: '',
    agentPreset: 'codex',
    agentRuntimeInstallConfirmed: false,
    agentAdvancedOpen: false,
    agentChecking: false,
    agentInstallingCLI: false,
    agentCheckResults: {},
    agentForm: {
      name: '',
      image: 'utbench-agent-base:latest',
      timeout_seconds: 600,
      network_disabled: false,
      env_from_host: '',
      compatible_models: [],
      compatible_languages: ['python', 'go', 'java', 'cpp'],
      output_globs: 'generated_test.py\ntest_*.py\n*_test.go\n*Test.java\n*test*.cpp',
      command: '',
    },
    // Skill 管理弹窗
    skillModalOpen: false,
    skillScanning: false,
    skillSaving: false,
    skillUploading: false,
    skillCommandRunning: false,
    skillInstallMode: 'scan',
    skillCommand: '',
    skillCommandOutput: '',
    skillUploadFileName: '',
    skillFormError: '',
    skillPackages: [],  // 扫描到的可用包列表
    skillForm: {
      skill_root: '',
      names: [],
    },
    // Skill 详情/对比
    skillDetailOpen: false,
    skillDetail: null,
    selectedSkillNames: [],
    skillCompareOpen: false,
    compareSkills: [],
    // 数据库管理
    dbTab: 'overview', // overview | runs | evaluation | assets | disk
    dbOverview: null,
    dbRuns: [],
    dbResults: [],
    dbArtifacts: [],
    dbFacets: { runs: [], models: [], languages: [], env_groups: [] },
    dbLoading: false,
    dbIngesting: false,
    dbReportGenerating: false,
    dbIngestRunID: '',
    dbFilters: { run_id:'', model:'', language:'' },
    dbArtifactFilters: { run_id:'', kind:'' },
    dbReportSelection: { run_ids: [], models: [], languages: [], score_eligible_only: false, merge_target: '', dedup_mode: 'merge' },
    lastDBReportResult: null,
    pinnedReportID: localStorage.getItem('utbench_pinned_report') || '',
    // 新增：各表数据
    dbGenerationRuns: [],
    dbGeneratedCases: [],
    dbPromptRenderings: [],
    dbEvaluationRuns: [],
    dbEvaluationStages: [],
    dbAssetSubjects: [],
    dbSubjectVersions: [],
    dbAssetGenerations: [],
    dbAssetEvaluations: [],
    dbReuseExplain: null,
    dbDatasetSamples: [],
    dbDatasetSnapshots: [],
    dbModelConfigs: [],
    dbPromptProfiles: [],
    dbEvaluationEnvs: [],
    dbScorePolicies: [],
    dbReports: [],
    dbRunArtifacts: [],
    dbExperiments: [],
    automations: [],
    automationRecentRuns: [],
    notificationDeliveries: [],
    notificationChannels: [],
    automationLoading: false,
    automationSaving: false,
    automationFormOpen: false,
    automationFormError: '',
    automationForm: {},
    automationPlan: { combinations: [], models: [], languages: [] },
    automationAdvancedOpen: false,
    channelFormOpen: false,
    channelFormError: '',
    channelForm: {},
    channelConfig: {},
    channelAdvancedOpen: false,
    channelTesting: {},
    channelTestResults: {},
    // 新增：筛选参数
    dbGenCaseFilter: { run_id:'', model:'', language:'' },
    dbEvalRunFilter: { run_id:'' },
    dbStageFilter: { evaluation_run_id:'' },
    dbSampleFilter: { language:'', class:'' },
    dbAssetFilter: { subject:'', language:'', sample:'' },
    dbReportFilter: { run_id:'' },
    dbRunArtifactFilter: { run_id:'' },
    form: {
      run_id:'', models:[], subjects:[], combinations:[{_id:1,framework:'model_api',model:'deepseek-v4-flash',skill:'no_skill'}], languages:[], class:'self_contained', scenario:'', level:'',
      max_samples:1, workers:4, mode:'full', phase:'full', source_run_id:'', manifest_path:'', evaluation_path:'',
      dry_run:false, reuse_generated:true, reuse_evaluation:false, mutation_enabled:true,
      mutation_timeout:360, mutation_policy:'warn', ingest:true, use_docker:true,
    },
    env: null,
    envChecking: false,
    toast: '',
    toastKind: 'ok', // ok | warn | err
    get toastStyle() {
      const map = {
        ok:  'background:var(--success-bg);color:var(--green)',
        warn:'background:var(--warn-bg);color:var(--yellow)',
        err: 'background:var(--error-bg);color:var(--red)',
      }
      return map[this.toastKind] || map.ok
    },
    get activeBuildImageName() {
      if (!this.env) return 'utbench:latest'
      if (this.buildTarget === 'agent') return this.env.agent_image_name || 'utbench-agent-base:latest'
      return this.env.eval_image_name || 'utbench:latest'
    },
    get activeBuildDockerfile() {
      if (this.buildTarget === 'agent') return 'docker/agents/Dockerfile'
      return 'Dockerfile'
    },
    buildModalOpen: false,
    buildTarget: 'eval',
    buildId: '',
    buildStatus: '',
    buildLogs: [],
    buildError: '',
    buildSSE: null,
    environment: null,
    environmentLoading: false,
    apiKeys: [],
    apiKeysLoading: false,
    installConfirmOpen: false,
    installTarget: null,
    installRunning: false,
    installOutput: '',
    installError: '',
    formSubmitting: false,
    formError: '',
    currentRun: null,
    currentLogs: [],
    currentReport: null,
    detailTab: 'logs',
    sseSource: null,
    runActionBusy: '',

    // 模型管理
    models: [],
    modelsLoading: false,
    modelFormOpen: false,
    modelFormMode: 'create', // create | edit
    modelForm: {
      name: '', enabled: true, provider: '', model_id: '',
      api_endpoint: '', anthropic_endpoint: '',
      api_key_env: '', api_key: '', api_key_set: false,
      parameters: { temperature: 0.7, top_p: 0.9, max_tokens: 4096 },
    },
    modelFormError: '',
    modelKeyVisible: {},
    theme: localStorage.getItem('utbench-theme') || 'dark',
    modelTesting: {},
    modelTestingAll: false,
    modelTestResults: {},
    _timerRuns: null,
    _timerDb: null,
    _timerEnv: null,
    _timerAutomation: null,
    _logCount: 0,
    _comboKey: 0,
    _comboSeq: 1,
    _comboAddLocked: false,
    partialsLoaded: false,

    async init() {
      this.applyTheme()
      await this.loadPartials()
      this.syncPageVisibility()
      await this.loadConfig()
      await this.loadModels()
      await this.loadEnv()
      await this.loadRuns()
      await this.loadDatabase()
      await this.loadAutomations()
      this._startTimers()
    },
    async loadPartials() {
      const nodes = Array.from(document.querySelectorAll('[data-partial]'))
      await Promise.all(nodes.map(async node => {
        const name = node.getAttribute('data-partial')
        if (!name) return
        const r = await fetch(`/partials/${name}.html`, { cache: 'no-store' })
        if (!r.ok) throw new Error(`load partial ${name}: HTTP ${r.status}`)
        if (window.Alpine?.destroyTree) {
          try { window.Alpine.destroyTree(node) } catch {}
        }
        node.innerHTML = await r.text()
        node._partialInitialized = false
      }))
      if (window.Alpine) {
        nodes.forEach(node => {
          if (node._partialInitialized) return
          Array.from(node.children).forEach(child => window.Alpine.initTree(child))
          node._partialInitialized = true
        })
      }
      this.partialsLoaded = true
      this.syncPageVisibility()
    },
    _startTimers() {
      this._stopTimers()
      this._timerRuns = setInterval(() => {
        if (this.page === 'dashboard' || this.page === 'runs' || this.page === 'run-detail') this.loadRuns()
      }, 4000)
      this._timerDb = setInterval(() => {
        if (this.page === 'database') this.loadDatabase()
      }, 12000)
      this._timerEnv = setInterval(() => {
        if (this.page === 'new-run' || this.page === 'models' || this.page === 'agents') this.loadEnv()
      }, 30000)
      this._timerAutomation = setInterval(() => {
        if (this.page === 'automations') this.loadAutomations()
      }, 10000)
    },
    _stopTimers() {
      if (this._timerRuns) { clearInterval(this._timerRuns); this._timerRuns = null }
      if (this._timerDb) { clearInterval(this._timerDb); this._timerDb = null }
      if (this._timerEnv) { clearInterval(this._timerEnv); this._timerEnv = null }
      if (this._timerAutomation) { clearInterval(this._timerAutomation); this._timerAutomation = null }
    },

    async loadConfig() {
      try { const r = await fetch('/api/config', { cache: 'no-store' }); this.config = await r.json(); this._comboKey++ }
      catch(e) { console.error('config', e) }
    },

    async loadEnv(force = false) {
      try {
        const r = await fetch(force ? '/api/env?refresh=1' : '/api/env')
        const prev = this.env
        this.env = await r.json()
        if (!prev && this.env.docker_available && this.env.eval_image_present) {
          const needDocker = this.env.os === 'windows' && !this.env.native_tools?.mutmut
          if (needDocker) this.form.use_docker = true
        }
      } catch(e) { console.error('env', e) }
    },

    async recheckEnv() {
      if (this.envChecking) return
      this.envChecking = true
      try {
        // 加短延迟让动画可见，同时 /api/env 每次都会重新探测 Docker/镜像/工具链
        const [r] = await Promise.all([fetch('/api/env?refresh=1'), new Promise(res => setTimeout(res, 400))])
        if (!r.ok) throw new Error('HTTP ' + r.status)
        this.env = await r.json()
        const parts = []
        parts.push(this.env.docker_available ? 'Docker✓' : 'Docker✗')
        parts.push(`评测镜像 ${this.env.eval_image_present ? '✓' : '✗'}`)
        parts.push(`拓扑 ${this.envTopologyText(this.env.topology_mode)}`)
        const nativeCount = Object.values(this.env.native_tools || {}).filter(Boolean).length
        const nativeTotal = Object.keys(this.env.native_tools || {}).length
        parts.push(`原生工具 ${nativeCount}/${nativeTotal}`)
        const kind = (!this.env.docker_available || !this.env.eval_image_present) ? 'warn' : 'ok'
        this.showToast('环境检查完成 · ' + parts.join(' · '), kind)
      } catch(e) {
        this.showToast('环境检查失败：' + e.message, 'err')
      } finally {
        this.envChecking = false
      }
    },

    async loadEnvironment(toast = false) {
      if (this.environmentLoading) return
      this.environmentLoading = true
      try {
        const r = await fetch('/api/environment/check', { cache: 'no-store' })
        const data = await r.json()
        if (!r.ok) throw new Error(data.error || 'HTTP ' + r.status)
        this.environment = data
        if (toast) {
          const s = data.summary || {}
          const kind = (s.missing || s.warning) ? 'warn' : 'ok'
          this.showToast(`环境检查完成：正常 ${s.ok || 0}，缺失 ${s.missing || 0}`, kind)
        }
      } catch (e) {
        this.showToast('环境检查失败：' + (e.message || String(e)), 'err')
      } finally {
        this.environmentLoading = false
      }
    },

    envStatusText(s) {
      return ({ ok:'正常', docker_ok:'Docker 可用', warning:'警告', missing:'缺失', unknown:'未验证' })[s] || s || '未知'
    },

    envStatusStyle(s) {
      const map = {
        ok: 'background:var(--success-bg);color:var(--green)',
        docker_ok: 'background:rgba(14,165,233,.1);color:var(--accent)',
        warning: 'background:var(--warn-bg);color:var(--yellow)',
        missing: 'background:var(--error-bg);color:var(--red)',
        unknown: 'background:var(--badge-bg);color:var(--fg-muted)',
      }
      return map[s] || map.unknown
    },

    // ─── API Key 管理 ──────────────────────────────────────────
    async loadAPIKeys() {
      this.apiKeysLoading = true
      try {
        const r = await fetch('/api/settings/api-keys', { cache: 'no-store' })
        const data = await r.json()
        if (!r.ok) throw new Error(data.error || 'HTTP ' + r.status)
        this.apiKeys = (data.keys || []).map(k => ({ ...k, _input: '' }))
      } catch (e) {
        this.showToast('加载 API Key 失败：' + (e.message || String(e)), 'err')
      } finally {
        this.apiKeysLoading = false
      }
    },

    async saveAPIKey(k) {
      if (!k._input || !k._input.trim()) return
      try {
        const body = { keys: { [k.key]: k._input.trim() } }
        const r = await fetch('/api/settings/api-keys', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(body),
        })
        const data = await r.json()
        if (!r.ok) throw new Error(data.error || 'HTTP ' + r.status)
        this.showToast(`已保存 ${k.label}`, 'ok')
        k._input = ''
        await this.loadAPIKeys()
      } catch (e) {
        this.showToast('保存失败：' + (e.message || String(e)), 'err')
      }
    },

    installCommandText(item) {
      const parts = item?.command_preview || []
      return parts.length ? parts.join(' ') : '由后端白名单生成'
    },

    openInstallConfirm(item) {
      this.installTarget = item
      this.installOutput = ''
      this.installError = ''
      this.installConfirmOpen = true
    },

    closeInstallConfirm() {
      if (this.installRunning) return
      this.installConfirmOpen = false
      this.installTarget = null
      this.installOutput = ''
      this.installError = ''
    },

    async confirmInstall() {
      if (!this.installTarget || this.installRunning) return
      this.installRunning = true
      this.installOutput = ''
      this.installError = ''
      try {
        const r = await fetch('/api/environment/install', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ tool: this.installTarget.id, confirmed: true }),
        })
        const data = await r.json()
        this.installOutput = data.output || ''
        if (!r.ok) {
          this.installError = data.error || 'HTTP ' + r.status
          throw new Error(this.installError)
        }
        this.showToast(`${this.installTarget.name} 安装完成`, 'ok')
        await this.loadEnvironment(false)
        await this.loadEnv()
        // 安装成功后延迟关闭弹窗（让用户看到输出）
        setTimeout(() => { this.closeInstallConfirm() }, 1500)
      } catch (e) {
        if (!this.installError) this.installError = e.message || String(e)
        this.showToast(`${this.installTarget?.name || '工具'} 安装失败：${this.installError}`, 'err', 6000)
      } finally {
        this.installRunning = false
      }
    },

    showToast(msg, kind = 'ok', ms = 3500) {
      this.toast = msg
      this.toastKind = kind
      clearTimeout(this._toastTimer)
      this._toastTimer = setTimeout(() => { this.toast = '' }, ms)
    },

    openBuildImage(target = 'eval') {
      const nextTarget = target || 'eval'
      if (this.buildTarget !== nextTarget) {
        this.buildId = ''
        this.buildStatus = ''
        this.buildLogs = []
        this.buildError = ''
      }
      this.buildTarget = nextTarget
      this.buildModalOpen = true
      this.reattachBuild(this.buildTarget)
    },
    closeBuildModal() { this.buildModalOpen = false; this.stopBuildSSE() },

    buildStatusColor(s) {
      const map = { running: 'color:var(--accent)', pending: 'color:var(--yellow)', completed: 'color:var(--green)', failed: 'color:var(--red)', canceled: 'color:var(--fg-muted)' }
      return map[s] || 'color:var(--fg-muted)'
    },

    defaultBuildTarget() {
      return 'eval'
    },

    buildTargetText(target) {
      const map = {
        eval: '评测镜像',
        evaluation: '评测镜像',
        agent: 'Agent 沙箱镜像',
        agents: 'Agent 沙箱镜像',
      }
      return map[target || 'eval'] || 'Docker 镜像'
    },

    envTopologyText(mode) {
      const map = {
        'container': '容器内运行',
        'host+docker': '宿主机 + Docker',
        'host': '宿主机',
      }
      return map[mode] || '未知'
    },

    envTopologyStyle(mode) {
      const map = {
        'container_control+nested_docker': 'background:rgba(16,185,129,.1);color:var(--green)',
        'container_control': 'background:rgba(245,158,11,.12);color:var(--yellow)',
        'host_control+docker_available': 'background:rgba(14,165,233,.1);color:var(--accent)',
        'host_control': 'background:rgba(100,116,139,.14);color:var(--fg-muted)',
      }
      return map[mode] || 'background:var(--badge-bg);color:var(--fg-muted)'
    },

    envImageStatusText(present) {
      return present ? '就绪' : '缺失'
    },

    envImageStatusStyle(present) {
      return present
        ? 'background:rgba(16,185,129,.1);color:var(--green)'
        : 'background:rgba(245,158,11,.12);color:var(--yellow)'
    },

    async reattachBuild(preferredTarget = '') {
      try {
        const r = await fetch('/api/env/build-image')
        const data = await r.json()
        if (data && data.build_id) {
          const active = data.status === 'running' || data.status === 'pending'
          if (preferredTarget && data.target && data.target !== preferredTarget && !active) return
          this.buildId = data.build_id
          this.buildTarget = data.target || this.buildTarget || 'eval'
          this.buildStatus = data.status
          this.buildError = data.error || ''
          const detail = await fetch('/api/env/build-image/' + this.buildId)
          if (detail.ok) { const d = await detail.json(); this.buildLogs = d.logs || [] }
          if (active) { this.startBuildSSE(this.buildId) }
        }
      } catch(e) { console.error('reattach build', e) }
    },

    async startBuild() {
      this.buildLogs = []; this.buildError = ''; this.buildStatus = 'pending'
      try {
        const r = await fetch('/api/env/build-image?target=' + encodeURIComponent(this.buildTarget || 'eval'), { method: 'POST' })
        const data = await r.json()
        if (!r.ok) { this.buildError = data.error || '启动失败'; this.buildStatus = 'failed'; return }
        this.buildId = data.build_id; this.buildStatus = data.status; this.buildTarget = data.target || this.buildTarget || 'eval'; this.startBuildSSE(this.buildId)
      } catch(e) { this.buildError = String(e); this.buildStatus = 'failed' }
    },

    async cancelBuild() {
      if (!this.buildId || (this.buildStatus !== 'running' && this.buildStatus !== 'pending')) return
      try {
        const r = await fetch('/api/env/build-image/' + this.buildId, { method: 'DELETE' })
        const data = await r.json()
        if (!r.ok) throw new Error(data.error || '取消失败')
        this.buildStatus = data.status || 'canceled'
        this.buildError = data.error || ''
        this.stopBuildSSE()
        this.showToast('Docker 构建已取消', 'warn')
      } catch(e) {
        this.showToast('取消构建失败：' + (e.message || String(e)), 'err', 5000)
      }
    },

    startBuildSSE(id) {
      this.stopBuildSSE()
      this.buildSSE = new EventSource('/api/env/build-image/' + id + '/events')
      this.buildSSE.onmessage = (e) => {
        try {
          const obj = JSON.parse(e.data)
          if (obj.type === 'log') { this.buildLogs.push(obj.payload); this.$nextTick(() => { const el = document.getElementById('build-log-bottom'); if(el) el.scrollIntoView({behavior:'smooth'}) }) }
          else if (obj.type === 'done') { this.buildStatus = obj.payload; this.stopBuildSSE(); this.loadEnv(true) }
        } catch {}
      }
      this.buildSSE.onerror = () => this.stopBuildSSE()
    },

    stopBuildSSE() { if (this.buildSSE) { this.buildSSE.close(); this.buildSSE = null } },

    async loadRuns() {
      this.runsLoading = true
      try {
        const r = await fetch('/api/runs')
        this.runs = await r.json()
        if (this.currentRun) {
          const updated = this.runs.find(r => r.run_id === this.currentRun.run_id)
          if (updated) { this.currentRun.status = updated.status; this.currentRun.ended_at = updated.ended_at; this.currentRun.error = updated.error }
        }
        this.autoSelectPinnedReport()
      } catch(e) { console.error('runs', e) }
      this.runsLoading = false
    },

    get filteredRuns() {
      return this.runs.filter(r => {
        const q = this.runFilter.toLowerCase()
        const subjectText = this.runSubjectSummary(r).toLowerCase()
        if (q && !r.run_id.includes(q) && !subjectText.includes(q) && !(r.spec?.models ?? []).join(',').toLowerCase().includes(q) && !(r.spec?.languages ?? []).join(',').toLowerCase().includes(q)) return false
        if (this.statusFilter && r.status !== this.statusFilter) return false
        return true
      })
    },

    // ─── 评测资产（磁盘扫描） ───────────────────────────────────────────────
    async loadAssets() {
      this.assetsLoading = true
      try {
        const r = await fetch('/api/assets/runs', { cache: 'no-store' })
        if (!r.ok) throw new Error('HTTP ' + r.status)
        this.assets = await r.json()
        // 清理已选中但已不存在的 ID
        const ids = new Set(this.assets.map(a => a.run_id))
        this.selectedAssetIds = this.selectedAssetIds.filter(id => ids.has(id))
      } catch (e) {
        console.error('[loadAssets]', e)
        this.showToast('加载资产失败：' + e.message, 'err')
      } finally {
        this.assetsLoading = false
      }
    },

    // 磁盘运行 Tab 的入口，复用 loadAssets 逻辑
    async loadDiskRuns() {
      await this.loadAssets()
    },

    get assetCountByQuality() {
      const out = { ok: 0, partial: 0, dirty: 0 }
      for (const a of this.assets) {
        if (a.quality_flag in out) out[a.quality_flag]++
      }
      return out
    },

    get assetAllLanguages() {
      const set = new Set()
      for (const a of this.assets) {
        for (const l of (a.languages || [])) set.add(l)
      }
      return Array.from(set).sort()
    },

    get assetFilteredRuns() {
      const q = (this.assetFilter || '').toLowerCase().trim()
      return this.assets.filter(a => {
        if (this.assetQualityFilter && a.quality_flag !== this.assetQualityFilter) return false
        if (this.assetLanguageFilter && !(a.languages || []).includes(this.assetLanguageFilter)) return false
        if (this.assetOnlyHasReport && !a.has_report_html) return false
        if (this.assetOnlyIngested && !a.ingested) return false
        if (q) {
          const hay = [
            a.run_id || '',
            a.label || '',
            (a.subjects || []).join(','),
            (a.models || []).join(','),
            (a.languages || []).join(','),
          ].join(' ').toLowerCase()
          if (!hay.includes(q)) return false
        }
        return true
      })
    },

    get allAssetsSelected() {
      const visible = this.assetFilteredRuns
      if (visible.length === 0) return false
      return visible.every(a => this.selectedAssetIds.includes(a.run_id))
    },

    toggleAssetSelection(runId) {
      const idx = this.selectedAssetIds.indexOf(runId)
      if (idx >= 0) this.selectedAssetIds.splice(idx, 1)
      else this.selectedAssetIds.push(runId)
    },

    toggleAllAssets(checked) {
      const visibleIds = this.assetFilteredRuns.map(a => a.run_id)
      if (checked) {
        const merged = new Set([...this.selectedAssetIds, ...visibleIds])
        this.selectedAssetIds = Array.from(merged)
      } else {
        const visibleSet = new Set(visibleIds)
        this.selectedAssetIds = this.selectedAssetIds.filter(id => !visibleSet.has(id))
      }
    },

    async deleteAsset(runId) {
      if (!confirm(`确定删除资产 ${runId}？\n会删除磁盘上该 run 的所有文件（生成、评测、报告），不可恢复。`)) return
      try {
        const r = await fetch('/api/runs/' + encodeURIComponent(runId), { method: 'DELETE' })
        const data = await r.json().catch(() => ({}))
        if (!r.ok) { this.showToast('删除失败：' + (data.error || 'HTTP ' + r.status), 'err'); return }
        this.showToast('已删除 ' + runId, 'ok')
        await this.loadAssets()
      } catch (e) { this.showToast('删除失败：' + e.message, 'err') }
    },

    async bulkDeleteAssets() {
      if (!this.selectedAssetIds.length) return
      const ids = [...this.selectedAssetIds]
      if (!confirm(`将批量删除 ${ids.length} 项资产，不可恢复。继续？`)) return
      let ok = 0, fail = 0
      for (const id of ids) {
        try {
          const r = await fetch('/api/runs/' + encodeURIComponent(id), { method: 'DELETE' })
          if (r.ok) ok++; else fail++
        } catch { fail++ }
      }
      this.selectedAssetIds = []
      this.showToast(`批量删除完成：成功 ${ok}，失败 ${fail}`, fail ? 'warn' : 'ok')
      await this.loadAssets()
    },

    async bulkIngestAssets() {
      if (!this.selectedAssetIds.length) return
      const ids = [...this.selectedAssetIds]
      if (!confirm(`将批量补录 ${ids.length} 项资产到数据库。继续？`)) return
      let ok = 0, fail = 0
      for (const id of ids) {
        try {
          await this.ingestRunToDB(id)
          ok++
        } catch { fail++ }
      }
      this.selectedAssetIds = []
      this.showToast(`批量补录完成：成功 ${ok}，失败 ${fail}`, fail ? 'warn' : 'ok')
      await this.loadDiskRuns()
    },

    formatBytes(n) {
      if (n == null || isNaN(n)) return '—'
      const u = ['B','KB','MB','GB','TB']
      let v = Number(n), i = 0
      while (v >= 1024 && i < u.length - 1) { v /= 1024; i++ }
      return v.toFixed(v < 10 && i > 0 ? 1 : 0) + ' ' + u[i]
    },

    get completedRuns() {
      // 返回已完成的任务，用于数据源选择
      return this.runs.filter(r => r.status === 'completed')
    },

    get mergedReports() {
      // 返回已完成的合并报告，用于增量合并选择
      return this.runs.filter(r => r.is_merged_report && r.status === 'completed')
    },

    get sortedMergedReports() {
      // 置顶的排最前面，其余按时间倒序
      const pinned = this.pinnedReportID
      return [...this.mergedReports].sort((a, b) => {
        if (a.run_id === pinned && b.run_id !== pinned) return -1
        if (b.run_id === pinned && a.run_id !== pinned) return 1
        return (b.started_at || '').localeCompare(a.started_at || '')
      })
    },

    get selectedScenarioCount() {
      if (this.form.phase !== 'full' && this.form.phase !== 'generate') return 0
      if (this.form.scenario) return 1
      return (this.config?.scenarios ?? []).length || 4
    },

    get selectedLanguageCount() {
      return this.form.languages.length
    },

    get selectedSubjectCount() {
      return (this.form.combinations || []).filter(c => c.model).length
    },

    get selectedModelCount() {
      return this.form.models.length
    },

    get selectedExecutionTargetCount() {
      return this.selectedSubjectCount || this.selectedModelCount
    },

    get selectedFrameworkCount() {
      return new Set((this.form.combinations || []).map(c => c.framework).filter(Boolean)).size
    },

    get selectedSkillCount() {
      return new Set((this.form.combinations || []).map(c => c.skill).filter(s => s && s !== 'no_skill')).size
    },

    get subjectPlanHint() {
      if (this.selectedSubjectCount > 0) {
        return `将按 ${this.selectedSubjectCount} 个 subject 下发任务；模型列表只用于补齐这些 subject 所引用的模型配置。`
      }
      if (this.selectedModelCount > 0) {
        return '当前未选择 subject，将按纯模型 baseline 执行。'
      }
      return '添加组合或退回到纯模型 baseline。'
    },

    get estimatedTaskCount() {
      const samples = Number(this.form.max_samples || 0)
      if (samples <= 0 || this.selectedExecutionTargetCount === 0 || this.selectedLanguageCount === 0) return '—'
      return this.selectedExecutionTargetCount * this.selectedLanguageCount * this.selectedScenarioCount * samples
    },

    get estimatedTaskFormula() {
      const samples = Number(this.form.max_samples || 0)
      if (samples <= 0) return '（样本无限制，实际数量由数据集决定）'
      return `（${this.selectedExecutionTargetCount} 被测对象 × ${this.selectedLanguageCount} 语言 × ${this.selectedScenarioCount} 场景 × ${samples} 样本）`
    },

    async loadDatabase() {
      this.dbLoading = true
      try {
        const [overview, runs, results, artifacts, facets] = await Promise.all([
          fetch('/api/db/overview?limit=8', { cache: 'no-store' }).then(r => r.json()),
          fetch('/api/db/runs?limit=50', { cache: 'no-store' }).then(r => r.json()),
          fetch('/api/db/results?limit=80', { cache: 'no-store' }).then(r => r.json()),
          fetch('/api/db/artifacts?limit=80', { cache: 'no-store' }).then(r => r.json()),
          fetch('/api/db/facets?limit=100', { cache: 'no-store' }).then(r => r.json()),
        ])
        this.dbOverview = overview
        this.dbRuns = Array.isArray(runs) ? runs : []
        this.dbResults = Array.isArray(results) ? results : []
        this.dbArtifacts = Array.isArray(artifacts) ? artifacts : []
        this.dbFacets = facets && !facets.error ? facets : { runs: [], models: [], languages: [], env_groups: [] }
      } catch(e) {
        this.showToast('加载数据库失败：' + (e.message || String(e)), 'err')
      } finally {
        this.dbLoading = false
      }
    },

    async loadAutomations() {
      this.automationLoading = true
      try {
        const [schedules, channels, deliveries] = await Promise.all([
          fetch('/api/automations', { cache: 'no-store' }).then(r => r.json()),
          fetch('/api/notification-channels', { cache: 'no-store' }).then(r => r.json()),
          fetch('/api/notification-deliveries?limit=60', { cache: 'no-store' }).then(r => r.json()),
        ])
        this.automations = Array.isArray(schedules) ? schedules : []
        this.notificationChannels = Array.isArray(channels) ? channels : []
        this.notificationDeliveries = Array.isArray(deliveries) ? deliveries : []
        const runs = await Promise.all(this.automations.map(a => fetch(`/api/automations/${a.schedule_id}/runs?limit=8`, { cache: 'no-store' }).then(r => r.json()).catch(() => [])))
        this.automationRecentRuns = runs.flat().filter(Boolean).sort((a, b) => String(b.created_at_utc || '').localeCompare(String(a.created_at_utc || ''))).slice(0, 30)
      } catch(e) {
        this.showToast('加载定时任务失败：' + (e.message || String(e)), 'err')
      } finally {
        this.automationLoading = false
      }
    },

    automationName(id) {
      return this.automations.find(a => a.schedule_id === id)?.name || id || '—'
    },

    defaultAutomationPlan() {
      const models = this._comboModels('model_api')
      const defaultModel = this.form.models?.[0] || (models.includes('deepseek-v4-flash') ? 'deepseek-v4-flash' : (models[0] || ''))
      return {
        mode: this.form.mode || 'full',
        phase: this.form.phase || 'full',
        source_run_id: this.form.source_run_id || '',
        manifest_path: this.form.manifest_path || '',
        evaluation_path: this.form.evaluation_path || '',
        combinations: [{ _id: 1, framework: 'model_api', model: defaultModel, skill: 'no_skill' }],
        models: [],
        languages: this.form.languages?.length ? [...this.form.languages] : ['python'],
        class: this.form.class || 'self_contained',
        scenario: this.form.scenario || '',
        level: this.form.level || '',
        max_samples: this.form.max_samples || 1,
        workers: this.form.workers || 4,
        dry_run: false,
        reuse_generated: true,
        reuse_evaluation: false,
        mutation_enabled: true,
        mutation_timeout: 1800,
        mutation_policy: 'warn',
        ingest: true,
      }
    },

    automationPlanFromSchedule(item) {
      let spec = {}
      let opts = {}
      try { spec = JSON.parse(item.run_spec_json || '{}') } catch {}
      try { opts = JSON.parse(item.orchestrator_options_json || '{}') } catch {}
      const subjects = Array.isArray(spec.subjects) ? spec.subjects : []
      const combos = subjects.map((subject, idx) => this.comboFromSubjectId(subject, idx + 1)).filter(Boolean)
      const models = Array.isArray(spec.models) ? spec.models : []
      return {
        mode: spec.mode || 'full',
        phase: opts.phase || 'full',
        source_run_id: opts.source_run_id || '',
        manifest_path: opts.manifest_path || '',
        evaluation_path: opts.evaluation_path || '',
        combinations: combos.length ? combos : [{ _id: 1, framework: 'model_api', model: models[0] || this._comboModels('model_api')[0] || '', skill: 'no_skill' }],
        models: combos.length ? [] : models,
        languages: Array.isArray(spec.languages) ? spec.languages : [],
        class: Array.isArray(spec.dataset_classes) ? spec.dataset_classes.join(',') : '',
        scenario: spec.dataset_scenario || '',
        level: spec.dataset_level || '',
        max_samples: Number(spec.max_samples || 0),
        workers: Number(spec.workers || 0),
        dry_run: !!spec.dry_run,
        reuse_generated: spec.reuse_generated !== false,
        reuse_evaluation: !!spec.reuse_evaluation,
        mutation_enabled: spec.mutation_enabled !== false,
        mutation_timeout: Number(spec.mutation_timeout_seconds || 1800),
        mutation_policy: spec.mutation_error_policy || 'warn',
        ingest: opts.ingest !== false,
      }
    },

    comboFromSubjectId(subject, id) {
      const parts = String(subject || '').split('__')
      if (parts.length < 3) return null
      return { _id: id, framework: parts[0] || 'model_api', model: parts[1] || '', skill: parts.slice(2).join('__') || 'no_skill' }
    },

    buildAutomationRunSpecFromPlan() {
      const p = this.automationPlan || {}
      const combos = Array.isArray(p.combinations) ? p.combinations : []
      const subjects = combos.filter(c => c.model).map(c => this.buildSubjectId(c))
      const models = subjects.length ? this.deriveModelsFromAutomationSubjects(subjects, combos) : (p.models || [])
      return {
        models,
        subjects,
        languages: p.languages || [],
        dataset_classes: String(p.class || '').split(',').map(s => s.trim()).filter(Boolean),
        dataset_scenario: p.scenario || '',
        dataset_level: p.level || '',
        mode: p.mode || 'full',
        dry_run: !!p.dry_run,
        reuse_generated: !!p.reuse_generated,
        reuse_evaluation: !!p.reuse_evaluation,
        mutation_enabled: !!p.mutation_enabled,
        mutation_timeout_seconds: Number(p.mutation_timeout || 1800),
        mutation_error_policy: p.mutation_policy || 'warn',
        max_samples: Number(p.max_samples || 0),
        workers: Number(p.workers || 0),
      }
    },

    buildAutomationOptionsFromPlan() {
      const p = this.automationPlan || {}
      return {
        phase: p.phase || 'full',
        source_run_id: p.source_run_id || '',
        manifest_path: p.manifest_path || '',
        evaluation_path: p.evaluation_path || '',
        ingest: p.ingest !== false,
      }
    },

    deriveModelsFromAutomationSubjects(subjects, combos) {
      const out = []
      const seen = new Set()
      for (const combo of combos || []) {
        if (!combo.model || seen.has(combo.model)) continue
        seen.add(combo.model); out.push(combo.model)
      }
      return out.length ? out : subjects.map(s => String(s).split('__')[1]).filter(Boolean)
    },

    refreshAutomationJSONFromPlan() {
      this.automationForm.run_spec_json = JSON.stringify(this.buildAutomationRunSpecFromPlan(), null, 2)
      this.automationForm.orchestrator_options_json = JSON.stringify(this.buildAutomationOptionsFromPlan(), null, 2)
    },

    get automationSubjectCount() {
      return (this.automationPlan.combinations || []).filter(c => c.model).length
    },

    get automationModelCount() {
      return (this.automationPlan.models || []).length
    },

    get automationExecutionTargetCount() {
      return this.automationSubjectCount || this.automationModelCount
    },

    get automationScenarioCount() {
      const p = this.automationPlan || {}
      if (p.phase !== 'full' && p.phase !== 'generate') return 0
      if (p.scenario) return 1
      return (this.config?.scenarios ?? []).length || 4
    },

    get automationEstimatedTaskCount() {
      const p = this.automationPlan || {}
      const samples = Number(p.max_samples || 0)
      const langs = (p.languages || []).length
      if (samples <= 0 || !this.automationExecutionTargetCount || !langs) return '—'
      return this.automationExecutionTargetCount * langs * this.automationScenarioCount * samples
    },

    get automationEstimatedTaskFormula() {
      const p = this.automationPlan || {}
      const samples = Number(p.max_samples || 0)
      if (samples <= 0) return '样本无限制，实际数量由数据集决定'
      return `${this.automationExecutionTargetCount} 被测对象 × ${(p.languages || []).length} 语言 × ${this.automationScenarioCount} 场景 × ${samples} 样本`
    },

    addAutomationCombination() {
      const models = this._comboModels('model_api')
      const defaultModel = models.includes('deepseek-v4-flash') ? 'deepseek-v4-flash' : (models[0] || '')
      if (!Array.isArray(this.automationPlan.combinations)) this.automationPlan.combinations = []
      this.automationPlan.models = []
      this.automationPlan.combinations.push({ _id: ++this._comboSeq, framework: 'model_api', model: defaultModel, skill: 'no_skill' })
    },

    removeAutomationCombination(idx) {
      if (!Array.isArray(this.automationPlan.combinations) || this.automationPlan.combinations.length <= 1) return
      this.automationPlan.combinations.splice(idx, 1)
    },

    onAutomationCombinationFrameworkChange(idx) {
      const combo = this.automationPlan.combinations?.[idx]
      if (!combo) return
      const models = this._comboModels(combo.framework)
      if (!models.includes(combo.model)) combo.model = models[0] || ''
      const skills = this._comboSkills(combo.framework)
      if (!skills.includes(combo.skill)) combo.skill = 'no_skill'
      this.automationPlan.models = []
    },

    onAutomationCombinationModelChange() {
      this.automationPlan.models = []
    },

    openAutomationForm() {
      this.automationFormError = ''
      this.automationPlan = this.defaultAutomationPlan()
      this.automationAdvancedOpen = false
      this.automationForm = {
        name: '',
        description: '',
        enabled: true,
        trigger_type: 'interval',
        cron_expr: '0 9 * * *',
        interval_seconds: 86400,
        timezone: 'Asia/Shanghai',
        concurrency_policy: 'skip',
        use_docker: true,
        run_spec_json: '',
        orchestrator_options_json: '',
        notify_policy_json: JSON.stringify({ on_success: true, on_failure: true, on_canceled: true, max_attempts: 3, channel_ids: [] }, null, 2),
      }
      this.refreshAutomationJSONFromPlan()
      this.automationFormOpen = true
    },

    editAutomation(item) {
      this.automationFormError = ''
      this.automationForm = { ...item }
      this.automationPlan = this.automationPlanFromSchedule(item)
      this.automationAdvancedOpen = false
      this.refreshAutomationJSONFromPlan()
      this.automationFormOpen = true
    },

    get automationPolicyChannels() {
      try {
        const policy = JSON.parse(this.automationForm.notify_policy_json || '{}')
        return policy.channel_ids || policy.channels || []
      }
      catch { return [] }
    },

    toggleAutomationPolicyChannel(id) {
      let policy
      try { policy = JSON.parse(this.automationForm.notify_policy_json || '{}') } catch { policy = {} }
      const list = Array.isArray(policy.channel_ids) ? policy.channel_ids : []
      policy.channel_ids = list.includes(id) ? list.filter(x => x !== id) : [...list, id]
      if (policy.on_success === undefined) policy.on_success = true
      if (policy.on_failure === undefined) policy.on_failure = true
      if (policy.on_canceled === undefined) policy.on_canceled = true
      if (policy.max_attempts === undefined) policy.max_attempts = 3
      this.automationForm.notify_policy_json = JSON.stringify(policy, null, 2)
    },

    async saveAutomation() {
      this.automationSaving = true
      this.automationFormError = ''
      try {
        const p = this.automationPlan || {}
        if (p.phase === 'generate' || p.phase === 'full') {
          if (!this.automationExecutionTargetCount) throw new Error('请至少选择一个 Subject 或模型')
          if (!(p.languages || []).length) throw new Error('请至少选择一种语言')
        }
        if (p.phase === 'evaluate' && !p.source_run_id && !p.manifest_path) {
          throw new Error('请选择来源任务或填写 manifest 路径')
        }
        if (p.phase === 'report' && !p.source_run_id && !p.evaluation_path) {
          throw new Error('请选择来源任务或填写评测结果路径')
        }
        this.refreshAutomationJSONFromPlan()
        JSON.parse(this.automationForm.run_spec_json || '{}')
        JSON.parse(this.automationForm.orchestrator_options_json || '{}')
        JSON.parse(this.automationForm.notify_policy_json || '{}')
        const id = this.automationForm.schedule_id
        const r = await fetch(id ? `/api/automations/${id}` : '/api/automations', {
          method: id ? 'PUT' : 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ ...this.automationForm, use_docker: true }),
        })
        const data = await r.json()
        if (!r.ok) throw new Error(data.error || 'HTTP ' + r.status)
        this.automationFormOpen = false
        this.showToast('定时计划已保存', 'ok')
        await this.loadAutomations()
      } catch(e) {
        this.automationFormError = e.message || String(e)
      } finally {
        this.automationSaving = false
      }
    },

    async deleteAutomation(id) {
      if (!confirm('删除这个定时计划？')) return
      const r = await fetch(`/api/automations/${id}`, { method: 'DELETE' })
      if (!r.ok) this.showToast('删除失败', 'err')
      await this.loadAutomations()
    },

    async triggerAutomation(id) {
      const r = await fetch(`/api/automations/${id}/trigger`, { method: 'POST' })
      const data = await r.json()
      if (!r.ok) { this.showToast(data.error || '触发失败', 'err'); return }
      this.showToast('已触发自动评测', 'ok')
      if (data.run_id) this.openRun(data.run_id)
      await this.loadAutomations()
    },

    channelGatewayTypes() {
      return [
        { type: 'webhook', label: '通用 Webhook', hint: 'POST 统一的任务 report.html 文件信息。' },
        { type: 'bot', label: '群机器人', hint: '飞书 / 钉钉 / 企业微信机器人，统一推送任务 report.html。' },
        { type: 'email_smtp', label: 'SMTP 邮件', hint: '通过企业邮箱或 SMTP 网关发送任务 report.html。' },
      ]
    },

    channelDisplayType(type = this.channelForm?.type) {
      return this.isBotChannelType(type) ? 'bot' : (type || 'webhook')
    },

    channelGatewayLabel(type) {
      const displayType = this.channelDisplayType(type)
      if (this.isBotChannelType(type)) return this.botProviderLabel(type)
      return this.channelGatewayTypes().find(x => x.type === displayType)?.label || type || '未知渠道'
    },

    isDefaultChannelName(name) {
      const value = String(name || '').trim()
      if (!value) return true
      return this.channelGatewayTypes().some(x => x.label === value) || ['飞书/Lark 机器人', '钉钉机器人', '企业微信机器人'].includes(value)
    },

    isBotChannelType(type) {
      return ['feishu', 'lark', 'dingtalk', 'wecom'].includes(type)
    },

    isGenericWebhookChannel() {
      return this.channelForm?.type === 'webhook'
    },

    isBotChannel() {
      return this.isBotChannelType(this.channelForm?.type)
    },

    isEmailChannel() {
      return this.channelForm?.type === 'email_smtp'
    },

    botProviderLabel(type) {
      const labels = { feishu: '飞书/Lark 机器人', lark: '飞书/Lark 机器人', dingtalk: '钉钉机器人', wecom: '企业微信机器人' }
      return labels[type] || '群机器人'
    },

    defaultChannelConfig(type = 'webhook') {
      if (type === 'email_smtp') {
        return {
          url: '',
          headers_text: '',
          auth_header: '',
          auth_env: '',
          secret: '',
          secret_env: '',
          host: 'smtp.example.com',
          port: 587,
          username: '',
          password: '',
          username_env: '',
          password_env: '',
          from: 'utbench@example.com',
          to_text: '',
          bot_provider: 'feishu',
        }
      }
      return {
        url: type === 'webhook' ? 'https://example.com/webhook' : '',
        headers_text: '',
        auth_header: '',
        auth_env: '',
        secret: '',
        secret_env: '',
        host: '',
        port: 587,
        username_env: '',
        password_env: '',
        username: '',
        password: '',
        from: '',
        to_text: '',
        bot_provider: this.isBotChannelType(type) ? type : 'feishu',
      }
    },

    parseHeaderLines(text) {
      const out = {}
      for (const line of String(text || '').split(/\r?\n/)) {
        const trimmed = line.trim()
        if (!trimmed) continue
        const idx = trimmed.indexOf(':')
        if (idx <= 0) continue
        out[trimmed.slice(0, idx).trim()] = trimmed.slice(idx + 1).trim()
      }
      return out
    },

    headerLinesFromObject(headers) {
      return Object.entries(headers || {}).map(([k, v]) => `${k}: ${v}`).join('\n')
    },

    listFromLines(text) {
      return String(text || '').split(/[\r\n,;]/).map(s => s.trim()).filter(Boolean)
    },

    channelConfigFromJSON(type, raw) {
      let cfg = {}
      try { cfg = JSON.parse(raw || '{}') } catch { cfg = {} }
      const next = this.defaultChannelConfig(type)
      next.url = cfg.url || ''
      next.headers_text = this.headerLinesFromObject(cfg.headers)
      next.auth_header = cfg.auth_header || ''
      next.auth_env = cfg.auth_env || ''
      next.bot_provider = this.isBotChannelType(type) ? type : (cfg.bot_provider || next.bot_provider || 'feishu')
      next.secret = cfg.secret || ''
      next.secret_env = cfg.secret_env || ''
      next.host = cfg.host || next.host
      next.port = Number(cfg.port || next.port || 587)
      next.username = cfg.username || (String(cfg.username_env || '').includes('@') ? cfg.username_env : '')
      next.password = cfg.password || (next.username && cfg.password_env ? cfg.password_env : '')
      next.username_env = next.username ? '' : (cfg.username_env || next.username_env)
      next.password_env = next.username ? '' : (cfg.password_env || next.password_env)
      next.from = cfg.from || next.from
      next.to_text = Array.isArray(cfg.to) ? cfg.to.join('\n') : ''
      return next
    },

    buildChannelConfigJSON() {
      const c = this.channelConfig || {}
      const type = this.channelForm.type
      if (type === 'email_smtp') {
        const cfg = {
          host: c.host || '',
          port: Number(c.port || 587),
          from: c.from || '',
          to: this.listFromLines(c.to_text),
        }
        if (c.username || c.password) {
          cfg.username = c.username || ''
          cfg.password = c.password || ''
        } else {
          if (c.username_env) cfg.username_env = c.username_env
          if (c.password_env) cfg.password_env = c.password_env
        }
        return JSON.stringify(cfg, null, 2)
      }
      const cfg = {
        url: c.url || '',
        headers: this.parseHeaderLines(c.headers_text),
      }
      if (type === 'webhook') {
        if (c.auth_header) cfg.auth_header = c.auth_header
        if (c.auth_env) cfg.auth_env = c.auth_env
      }
      if (type === 'feishu' || type === 'lark' || type === 'dingtalk') {
        if (c.secret) cfg.secret = c.secret
        if (c.secret_env) cfg.secret_env = c.secret_env
      }
      return JSON.stringify(cfg, null, 2)
    },

    refreshChannelConfigJSON() {
      this.channelForm.config_json = this.buildChannelConfigJSON()
    },

    setChannelType(type) {
      const prevName = this.channelForm.name
      const nextType = type === 'bot' ? (this.channelConfig?.bot_provider || 'feishu') : type
      this.channelForm.type = nextType
      this.channelConfig = this.defaultChannelConfig(nextType)
      if (this.isDefaultChannelName(prevName)) this.channelForm.name = this.channelGatewayLabel(nextType)
      this.refreshChannelConfigJSON()
    },

    setBotProvider(provider) {
      const prevName = this.channelForm.name
      const previous = this.channelConfig || {}
      this.channelForm.type = provider
      this.channelConfig = { ...this.defaultChannelConfig(provider), ...previous, bot_provider: provider }
      if (this.isDefaultChannelName(prevName)) this.channelForm.name = this.channelGatewayLabel(provider)
      this.refreshChannelConfigJSON()
    },

    toggleChannelAdvanced() {
      this.refreshChannelConfigJSON()
      this.channelAdvancedOpen = !this.channelAdvancedOpen
    },

    openChannelForm() {
      this.channelFormError = ''
      this.channelAdvancedOpen = false
      this.channelForm = { name: '', type: 'webhook', enabled: true, config_json: '{}' }
      this.channelConfig = this.defaultChannelConfig(this.channelForm.type)
      this.channelForm.name = this.channelGatewayLabel(this.channelForm.type)
      this.refreshChannelConfigJSON()
      this.channelFormOpen = true
    },

    async editChannel(channel) {
      this.channelFormError = ''
      this.channelAdvancedOpen = false
      let detail = channel
      if (channel?.channel_id) {
        try {
          const r = await fetch(`/api/notification-channels/${encodeURIComponent(channel.channel_id)}`, { cache: 'no-store' })
          if (r.ok) detail = await r.json()
        } catch {}
      }
      this.channelForm = { ...detail }
      if (this.channelForm.type === 'openclaw_gateway' || this.channelForm.type === 'notification_gateway') {
        this.channelForm.type = 'webhook'
      }
      if (!this.channelForm.config_json) this.channelForm.config_json = '{}'
      this.channelConfig = this.channelConfigFromJSON(this.channelForm.type, this.channelForm.config_json)
      this.refreshChannelConfigJSON()
      this.channelFormOpen = true
    },

    async saveChannel() {
      this.channelFormError = ''
      try {
        this.refreshChannelConfigJSON()
        const cfg = JSON.parse(this.channelForm.config_json || '{}')
        if (!String(this.channelForm.name || '').trim()) throw new Error('请填写渠道名称')
        if (this.isEmailChannel()) {
          if (!cfg.host || !cfg.from || !Array.isArray(cfg.to) || !cfg.to.length) throw new Error('请填写 SMTP host、发件人和收件人')
        } else if (this.isBotChannel()) {
          if (!cfg.url) throw new Error('请填写机器人 Webhook URL')
        } else if (!cfg.url) {
          throw new Error('请填写 Webhook / 网关 URL')
        }
        const id = this.channelForm.channel_id
        const r = await fetch(id ? `/api/notification-channels/${id}` : '/api/notification-channels', {
          method: id ? 'PUT' : 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(this.channelForm),
        })
        const data = await r.json()
        if (!r.ok) throw new Error(data.error || 'HTTP ' + r.status)
        this.channelFormOpen = false
        this.showToast('通知渠道已保存', 'ok')
        await this.loadAutomations()
      } catch(e) {
        this.channelFormError = e.message || String(e)
      }
    },

    isChannelTesting(id) {
      return !!id && this.channelTesting?.[id] === true
    },

    async testChannel(id) {
      if (!id || this.isChannelTesting(id)) return
      this.channelTesting = { ...this.channelTesting, [id]: true }
      this.channelTestResults = { ...this.channelTestResults, [id]: null }
      this.showToast('正在测试通知渠道...', 'warn')
      const ctrl = new AbortController()
      const timer = setTimeout(() => ctrl.abort(), 180000)
      try {
        const r = await fetch(`/api/notification-channels/${encodeURIComponent(id)}/test`, { method: 'POST', signal: ctrl.signal })
        const data = await r.json().catch(() => ({}))
        if (!r.ok) throw new Error(data.error || '通知测试失败')
        this.channelTestResults = { ...this.channelTestResults, [id]: { ok: true, message: data.response_summary || '连接测试成功', at: new Date().toLocaleTimeString() } }
        this.showToast('通知测试成功', 'ok')
        await this.loadAutomations()
      } catch(e) {
        const message = e.name === 'AbortError' ? '测试超时，请检查网络、SMTP/Webhook 地址或报告附件大小' : (e.message || String(e))
        this.channelTestResults = { ...this.channelTestResults, [id]: { ok: false, message, at: new Date().toLocaleTimeString() } }
        this.showToast('通知测试失败：' + message, 'err', 6000)
      } finally {
        clearTimeout(timer)
        this.channelTesting = { ...this.channelTesting, [id]: false }
      }
    },

    async deleteChannel(id) {
      if (!id || !confirm('删除这个通知渠道？')) return
      const r = await fetch(`/api/notification-channels/${encodeURIComponent(id)}`, { method: 'DELETE' })
      const data = await r.json().catch(() => ({}))
      if (!r.ok) { this.showToast(data.error || '删除通知渠道失败', 'err'); return }
      if (this.channelForm?.channel_id === id) this.channelFormOpen = false
      this.showToast('通知渠道已删除', 'ok')
      await this.loadAutomations()
    },

    async loadDBResults() {
      const q = new URLSearchParams()
      if (this.dbFilters.run_id) q.set('run_id', this.dbFilters.run_id)
      if (this.dbFilters.model) q.set('model', this.dbFilters.model)
      if (this.dbFilters.language) q.set('language', this.dbFilters.language)
      q.set('limit', '120')
      try {
        const r = await fetch('/api/db/results?' + q.toString(), { cache: 'no-store' })
        this.dbResults = await r.json()
      } catch(e) { this.showToast('加载结果失败：' + (e.message || String(e)), 'err') }
    },

    async loadDBArtifacts() {
      const q = new URLSearchParams()
      if (this.dbArtifactFilters.run_id) q.set('run_id', this.dbArtifactFilters.run_id)
      if (this.dbArtifactFilters.kind) q.set('kind', this.dbArtifactFilters.kind)
      q.set('limit', '120')
      try {
        const r = await fetch('/api/db/artifacts?' + q.toString(), { cache: 'no-store' })
        this.dbArtifacts = await r.json()
      } catch(e) { this.showToast('加载 artifact 失败：' + (e.message || String(e)), 'err') }
    },

    selectDBRun(runID) {
      this.dbFilters.run_id = runID
      this.dbArtifactFilters.run_id = runID
      if (!this.dbReportSelection.run_ids.includes(runID)) {
        this.dbReportSelection.run_ids = [runID]
      }
      this.loadDBResults()
      this.loadDBArtifacts()
    },

    async ingestRunToDB(runID = '') {
      const targetRunID = runID || this.dbIngestRunID
      if (!targetRunID || this.dbIngesting) return
      this.dbIngesting = true
      try {
        const r = await fetch('/api/db/ingest-run', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ run_id: targetRunID }),
        })
        const data = await r.json()
        if (!r.ok) throw new Error(data.error || 'HTTP ' + r.status)
        this.showToast(`已入库 ${data.run_id} · 结果 ${data.evaluation_results || 0} 条`, 'ok')
        if (!runID) this.dbIngestRunID = ''
        await this.loadDatabase()
      } catch(e) {
        this.showToast('入库失败：' + (e.message || String(e)), 'err', 6000)
      } finally {
        this.dbIngesting = false
      }
    },

    async generateDBReport() {
      if (this.dbReportGenerating) return
      this.dbReportGenerating = true
      this.lastDBReportResult = null
      try {
        const body = {
          source_run_ids: this.dbReportSelection.run_ids.length ? this.dbReportSelection.run_ids : (this.dbFilters.run_id ? [this.dbFilters.run_id] : []),
          models: this.dbReportSelection.models.length ? this.dbReportSelection.models : (this.dbFilters.model ? [this.dbFilters.model] : []),
          languages: this.dbReportSelection.languages.length ? this.dbReportSelection.languages : (this.dbFilters.language ? [this.dbFilters.language] : []),
          score_eligible_only: !!this.dbReportSelection.score_eligible_only,
          dedup_mode: this.dbReportSelection.dedup_mode || 'merge',
        }
        // 增量合并：如果选择了合并目标，传入 run_id 让后端追加到已有报告
        if (this.dbReportSelection.merge_target) {
          body.run_id = this.dbReportSelection.merge_target
        }
        const r = await fetch('/api/db/report', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(body),
        })
        const data = await r.json()
        if (!r.ok) throw new Error(data.error || 'HTTP ' + r.status)
        const action = this.dbReportSelection.merge_target ? '追加合并' : '生成'
        this.showToast(`DB 报告已${action}：${data.result_count || 0} 条结果`, 'ok')
        this.lastDBReportResult = data
        await this.loadRuns()
        await this.loadDatabase()
      } catch(e) {
        this.showToast('生成 DB 报告失败：' + (e.message || String(e)), 'err', 6000)
      } finally {
        this.dbReportGenerating = false
      }
    },

    selectMergeTarget(runID) {
      this.dbReportSelection.merge_target = runID
    },

    togglePinReport(runID) {
      if (this.pinnedReportID === runID) {
        // 取消置顶
        this.pinnedReportID = ''
        localStorage.removeItem('utbench_pinned_report')
        this.dbReportSelection.merge_target = ''
      } else {
        // 置顶并自动选为合并目标
        this.pinnedReportID = runID
        localStorage.setItem('utbench_pinned_report', runID)
        this.dbReportSelection.merge_target = runID
      }
    },

    autoSelectPinnedReport() {
      // 如果有置顶报告且存在于列表中，自动选为合并目标
      if (this.pinnedReportID && this.mergedReports.some(r => r.run_id === this.pinnedReportID)) {
        this.dbReportSelection.merge_target = this.pinnedReportID
      } else if (this.pinnedReportID) {
        // 置顶报告已不存在，清除
        this.pinnedReportID = ''
        localStorage.removeItem('utbench_pinned_report')
      }
    },

    toggleReportSelection(key, value) {
      const list = this.dbReportSelection[key] || []
      if (list.includes(value)) {
        this.dbReportSelection[key] = list.filter(v => v !== value)
      } else {
        this.dbReportSelection[key] = [...list, value]
      }
    },

    clearReportSelection() {
      this.dbReportSelection = { run_ids: [], models: [], languages: [], score_eligible_only: false, merge_target: '', dedup_mode: 'merge' }
      this.lastDBReportResult = null
      // 清空后仍保留置顶报告的自动选择
      this.$nextTick(() => this.autoSelectPinnedReport())
    },

    // 新增：各表加载函数
    async loadDBGenerationRuns() {
      try {
        const r = await fetch('/api/db/generation-runs?limit=100', { cache: 'no-store' })
        this.dbGenerationRuns = await r.json()
      } catch(e) { this.showToast('加载生成运行失败：' + e.message, 'err') }
    },

    async loadDBGeneratedCases() {
      const q = new URLSearchParams()
      if (this.dbGenCaseFilter.run_id) q.set('run_id', this.dbGenCaseFilter.run_id)
      if (this.dbGenCaseFilter.model) q.set('model', this.dbGenCaseFilter.model)
      if (this.dbGenCaseFilter.language) q.set('language', this.dbGenCaseFilter.language)
      q.set('limit', '100')
      try {
        const r = await fetch('/api/db/generated-cases?' + q.toString(), { cache: 'no-store' })
        this.dbGeneratedCases = await r.json()
      } catch(e) { this.showToast('加载生成样本失败：' + e.message, 'err') }
    },

    async loadDBPromptRenderings() {
      const q = new URLSearchParams()
      if (this.dbGenCaseFilter.run_id) q.set('run_id', this.dbGenCaseFilter.run_id)
      q.set('limit', '100')
      try {
        const r = await fetch('/api/db/prompt-renderings?' + q.toString(), { cache: 'no-store' })
        this.dbPromptRenderings = await r.json()
      } catch(e) { this.showToast('加载Prompt记录失败：' + e.message, 'err') }
    },

    async loadDBEvaluationRuns() {
      const q = new URLSearchParams()
      if (this.dbEvalRunFilter.run_id) q.set('run_id', this.dbEvalRunFilter.run_id)
      q.set('limit', '100')
      try {
        const r = await fetch('/api/db/evaluation-runs?' + q.toString(), { cache: 'no-store' })
        this.dbEvaluationRuns = await r.json()
      } catch(e) { this.showToast('加载评测运行失败：' + e.message, 'err') }
    },

    async loadDBEvaluationStages() {
      const q = new URLSearchParams()
      if (this.dbStageFilter.evaluation_run_id) q.set('evaluation_run_id', this.dbStageFilter.evaluation_run_id)
      q.set('limit', '100')
      try {
        const r = await fetch('/api/db/evaluation-stages?' + q.toString(), { cache: 'no-store' })
        this.dbEvaluationStages = await r.json()
      } catch(e) { this.showToast('加载Stage记录失败：' + e.message, 'err') }
    },

    async loadDBDatasetSamples() {
      const q = new URLSearchParams()
      if (this.dbSampleFilter.language) q.set('language', this.dbSampleFilter.language)
      if (this.dbSampleFilter.class) q.set('class', this.dbSampleFilter.class)
      q.set('limit', '100')
      try {
        const r = await fetch('/api/db/dataset-samples?' + q.toString(), { cache: 'no-store' })
        this.dbDatasetSamples = await r.json()
      } catch(e) { this.showToast('加载样本记录失败：' + e.message, 'err') }
    },

    async loadDBDatasetSnapshots() {
      try {
        const r = await fetch('/api/db/dataset-snapshots?limit=50', { cache: 'no-store' })
        this.dbDatasetSnapshots = await r.json()
      } catch(e) { this.showToast('加载快照记录失败：' + e.message, 'err') }
    },

    async loadDBAssetSubjects() {
      try {
        const r = await fetch('/api/db/asset-subjects?limit=500', { cache: 'no-store' })
        this.dbAssetSubjects = await r.json()
      } catch(e) { this.showToast('加载 subject 资产失败：' + e.message, 'err') }
    },

    async loadDBSubjectVersions() {
      const q = new URLSearchParams()
      if (this.dbAssetFilter.subject) q.set('subject', this.dbAssetFilter.subject)
      q.set('limit', '500')
      try {
        const r = await fetch('/api/db/subject-versions?' + q.toString(), { cache: 'no-store' })
        this.dbSubjectVersions = await r.json()
      } catch(e) { this.showToast('加载 subject version 失败：' + e.message, 'err') }
    },

    async loadDBAssetGenerations() {
      const q = new URLSearchParams()
      if (this.dbAssetFilter.subject) q.set('subject', this.dbAssetFilter.subject)
      if (this.dbAssetFilter.language) q.set('language', this.dbAssetFilter.language)
      if (this.dbAssetFilter.sample) q.set('sample', this.dbAssetFilter.sample)
      q.set('limit', '2000')
      try {
        const r = await fetch('/api/db/asset-generations?' + q.toString(), { cache: 'no-store' })
        this.dbAssetGenerations = await r.json()
      } catch(e) { this.showToast('加载生成资产失败：' + e.message, 'err') }
    },

    async loadDBAssetEvaluations() {
      const q = new URLSearchParams()
      if (this.dbAssetFilter.subject) q.set('subject', this.dbAssetFilter.subject)
      if (this.dbAssetFilter.language) q.set('language', this.dbAssetFilter.language)
      if (this.dbAssetFilter.sample) q.set('sample', this.dbAssetFilter.sample)
      q.set('limit', '2000')
      try {
        const r = await fetch('/api/db/asset-evaluations?' + q.toString(), { cache: 'no-store' })
        this.dbAssetEvaluations = await r.json()
      } catch(e) { this.showToast('加载评测资产失败：' + e.message, 'err') }
    },

    async loadDBReuseExplain() {
      this.dbReuseExplain = null
      if (!this.dbAssetFilter.subject || !this.dbAssetFilter.language || !this.dbAssetFilter.sample) return
      const q = new URLSearchParams()
      q.set('subject', this.dbAssetFilter.subject)
      q.set('language', this.dbAssetFilter.language)
      q.set('sample', this.dbAssetFilter.sample)
      q.set('limit', '20')
      try {
        const r = await fetch('/api/db/asset-explain-reuse?' + q.toString(), { cache: 'no-store' })
        this.dbReuseExplain = await r.json()
      } catch(e) { this.showToast('加载复用解释失败：' + e.message, 'err') }
    },

    async applyDBAssetFilter(subject = this.dbAssetFilter.subject, language = this.dbAssetFilter.language, sample = this.dbAssetFilter.sample) {
      this.dbAssetFilter.subject = subject || ''
      this.dbAssetFilter.language = language || ''
      this.dbAssetFilter.sample = sample || ''
      await Promise.all([
        this.loadDBSubjectVersions(),
        this.loadDBAssetGenerations(),
        this.loadDBAssetEvaluations(),
      ])
      await this.loadDBReuseExplain()
    },

    async loadDBModelConfigs() {
      try {
        const r = await fetch('/api/db/model-configs?limit=50', { cache: 'no-store' })
        this.dbModelConfigs = await r.json()
      } catch(e) { this.showToast('加载模型配置失败：' + e.message, 'err') }
    },

    async loadDBPromptProfiles() {
      try {
        const r = await fetch('/api/db/prompt-profiles?limit=20', { cache: 'no-store' })
        this.dbPromptProfiles = await r.json()
      } catch(e) { this.showToast('加载Prompt策略失败：' + e.message, 'err') }
    },

    async loadDBEvaluationEnvs() {
      try {
        const r = await fetch('/api/db/evaluation-envs?limit=20', { cache: 'no-store' })
        this.dbEvaluationEnvs = await r.json()
      } catch(e) { this.showToast('加载评测环境失败：' + e.message, 'err') }
    },

    async loadDBScorePolicies() {
      try {
        const r = await fetch('/api/db/score-policies?limit=10', { cache: 'no-store' })
        this.dbScorePolicies = await r.json()
      } catch(e) { this.showToast('加载评分策略失败：' + e.message, 'err') }
    },

    async loadDBReports() {
      const q = new URLSearchParams()
      if (this.dbReportFilter.run_id) q.set('run_id', this.dbReportFilter.run_id)
      q.set('limit', '50')
      try {
        const r = await fetch('/api/db/reports?' + q.toString(), { cache: 'no-store' })
        this.dbReports = await r.json()
      } catch(e) { this.showToast('加载报告记录失败：' + e.message, 'err') }
    },

    async loadDBRunArtifacts() {
      const q = new URLSearchParams()
      if (this.dbRunArtifactFilter.run_id) q.set('run_id', this.dbRunArtifactFilter.run_id)
      q.set('limit', '100')
      try {
        const r = await fetch('/api/db/run-artifacts?' + q.toString(), { cache: 'no-store' })
        this.dbRunArtifacts = await r.json()
      } catch(e) { this.showToast('加载关联记录失败：' + e.message, 'err') }
    },

    async loadDBExperiments() {
      try {
        const r = await fetch('/api/db/experiments?limit=20', { cache: 'no-store' })
        this.dbExperiments = await r.json()
      } catch(e) { this.showToast('加载实验记录失败：' + e.message, 'err') }
    },

    // 按当前Tab加载对应数据（仅首次进入时加载，避免重复查询）
    _dbTabLoaded: {},
    async loadDBTabData(force) {
      const tab = this.dbTab
      if (!force && this._dbTabLoaded[tab]) return
      this._dbTabLoaded[tab] = true
      if (tab === 'overview') await this.loadDatabase()
      else if (tab === 'runs') {
        await Promise.all([this.loadDBGenerationRuns(), this.loadDBGeneratedCases(), this.loadDBPromptRenderings()])
      }
      else if (tab === 'evaluation') {
        await Promise.all([this.loadDBEvaluationRuns(), this.loadDBEvaluationStages()])
        await this.loadDBResults()
      }
      else if (tab === 'assets') {
        // 分批加载：先加载核心数据，再按需加载次要数据
        await Promise.all([
          this.loadDBAssetSubjects(),
          this.loadDBAssetGenerations(),
          this.loadDBAssetEvaluations(),
        ])
        // 次要数据延迟加载，不阻塞 UI
        Promise.all([
          this.loadDBSubjectVersions(),
          this.loadDBDatasetSamples(),
          this.loadDBDatasetSnapshots(),
          this.loadDBModelConfigs(),
          this.loadDBPromptProfiles(),
          this.loadDBEvaluationEnvs(),
          this.loadDBScorePolicies(),
          this.loadDBReports(),
          this.loadDBArtifacts(),
          this.loadDBRunArtifacts(),
        ]).then(() => this.loadDBReuseExplain())
      }
      else if (tab === 'disk') {
        await this.loadDiskRuns()
      }
    },

    get dbReportSelectionSummary() {
      const s = this.dbReportSelection
      const runs = s.run_ids.length || '全部'
      const models = s.models.length || '全部'
      const langs = s.languages.length || '全部'
      return `运行 ${runs} · 模型 ${models} · 语言 ${langs}`
    },

    get dbAssetTree() {
      const subjectMeta = new Map((this.dbAssetSubjects || []).map(row => [row.subject_id, row]))
      const subjects = new Map()
      const ensureSubject = (subjectId) => {
        if (!subjects.has(subjectId)) {
          const meta = subjectMeta.get(subjectId) || {}
          subjects.set(subjectId, {
            subject_id: subjectId,
            subject_kind: meta.subject_kind || '',
            framework: meta.framework || '',
            model: meta.model || '',
            skill: meta.skill || '',
            generated_cases: meta.generated_cases || 0,
            evaluation_results: meta.evaluation_results || 0,
            latest_generated_at: meta.latest_generated_at || '',
            languages: new Map(),
          })
        }
        return subjects.get(subjectId)
      }
      const ensureLanguage = (subjectNode, language) => {
        if (!subjectNode.languages.has(language)) {
          subjectNode.languages.set(language, { language, samples: new Map() })
        }
        return subjectNode.languages.get(language)
      }
      const ensureSample = (langNode, sampleId) => {
        if (!langNode.samples.has(sampleId)) {
          langNode.samples.set(sampleId, {
            sample_id: sampleId,
            generations: [],
            evaluations: [],
            latest_generation: null,
            latest_evaluation: null,
          })
        }
        return langNode.samples.get(sampleId)
      }
      for (const row of (this.dbAssetGenerations || [])) {
        const subjectNode = ensureSubject(row.subject_id || 'unknown')
        const langNode = ensureLanguage(subjectNode, row.language || 'unknown')
        const sampleNode = ensureSample(langNode, row.sample_id || 'unknown')
        sampleNode.generations.push(row)
      }
      for (const row of (this.dbAssetEvaluations || [])) {
        const subjectNode = ensureSubject(row.subject_id || 'unknown')
        const langNode = ensureLanguage(subjectNode, row.language || 'unknown')
        const sampleNode = ensureSample(langNode, row.sample_id || 'unknown')
        sampleNode.evaluations.push(row)
      }
      const timeValue = (value) => value ? (Date.parse(value) || 0) : 0
      const out = Array.from(subjects.values()).map(subject => {
        subject.languages = Array.from(subject.languages.values()).map(lang => {
          lang.samples = Array.from(lang.samples.values()).map(sample => {
            sample.generations.sort((a, b) => timeValue(b.generated_at_utc) - timeValue(a.generated_at_utc))
            sample.evaluations.sort((a, b) => timeValue(b.created_at_utc) - timeValue(a.created_at_utc))
            sample.latest_generation = sample.generations[0] || null
            sample.latest_evaluation = sample.evaluations[0] || null
            return sample
          }).sort((a, b) => a.sample_id.localeCompare(b.sample_id))
          return lang
        }).sort((a, b) => a.language.localeCompare(b.language))
        return subject
      })
      out.sort((a, b) => a.subject_id.localeCompare(b.subject_id))
      return out
    },

    get selectedDBReportRuns() {
      const selected = this.dbReportSelection.run_ids
      if (!selected.length) return this.dbFacets.runs || []
      return (this.dbFacets.runs || []).filter(r => selected.includes(r.run_id))
    },

    get dbEnvWarning() {
      const runs = this.selectedDBReportRuns
      if (!runs.length) return ''
      const envs = [...new Set(runs.map(r => r.env_id || 'unknown'))]
      if (envs.length > 1) return `包含 ${envs.length} 个不同环境，报告仅供参考`
      if (envs[0] === 'unknown' || envs[0]?.startsWith('evaluation_env_unknown')) return '环境指纹未采集，报告仅供参考'
      return ''
    },

    get dbEnvCompatibilityStyle() {
      const warning = this.dbEnvWarning
      if (!warning) return 'background:var(--success-bg);color:var(--green)'
      if (warning.includes('不同环境')) return 'background:var(--warn-bg);color:var(--yellow)'
      return 'background:var(--error-bg);color:var(--red)'
    },

    get dbEnvCompatibilityText() {
      const warning = this.dbEnvWarning
      if (!warning) return '环境一致性：所选运行环境相同，可作为正式排名参考'
      return '环境检查：' + warning
    },

    get dbCards() {
      const d = this.dbOverview || {}
      return [
        { label:'生成运行', value:d.generation_runs ?? 0, color:'kpi-total' },
        { label:'评测运行', value:d.evaluation_runs ?? 0, color:'kpi-running' },
        { label:'生成样本', value:d.generated_cases ?? 0, color:'kpi-completed' },
        { label:'评测结果', value:d.evaluation_results ?? 0, color:'kpi-completed' },
        { label:'Artifacts', value:d.artifacts ?? 0, color:'kpi-total' },
        { label:'报告', value:d.reports ?? 0, color:'kpi-running' },
      ]
    },

    get dbHasData() {
      const d = this.dbOverview || {}
      return ['generation_runs', 'evaluation_runs', 'generated_cases', 'evaluation_results', 'artifacts', 'reports']
        .some(k => Number(d[k] || 0) > 0)
    },

    get dashStats() {
      const total = this.runs.length
      const running = this.runs.filter(r => r.status === 'running').length
      const completed = this.runs.filter(r => r.status === 'completed').length
      const failed = this.runs.filter(r => r.status === 'failed').length
      const canceled = this.runs.filter(r => r.status === 'canceled').length
      return [
        { label:'总任务数', value: total,   color: 'kpi-total' },
        { label:'运行中',    value: running, color: 'kpi-running' },
        { label:'已完成',    value: completed,color:'kpi-completed' },
        { label:'失败/取消', value: failed + canceled, color: 'kpi-failed' },
      ]
    },

    get submitButtonText() {
      const phaseLabels = {
        full: '开始评测',
        generate: '开始生成',
        evaluate: '开始评测',
        report: '生成报告',
      }
      return phaseLabels[this.form.phase] || '开始评测'
    },

    get submitDisabled() {
      if (this.form.phase === 'generate' || this.form.phase === 'full') {
        return this.selectedExecutionTargetCount === 0 || this.form.languages.length === 0
      }
      // evaluate 和 report 需要有数据源（source_run_id 或自定义路径）
      if (this.form.phase === 'evaluate') {
        return !this.form.source_run_id && !this.form.manifest_path
      }
      if (this.form.phase === 'report') {
        return !this.form.source_run_id && !this.form.evaluation_path
      }
      return false
    },

    statusText(s, paused) {
      const map = { pending:'等待中', running: paused ? '已暂停' : '运行中', completed:'已完成', failed:'失败', canceled:'已取消' }
      return (map[s] || s)
    },

    goto(id) {
      this.stopSSE(); this.page = id; this.syncPageVisibility()
      if (id === 'models') this.loadModels()
      if (id === 'agents') { this.loadAPIKeys(); this.loadModels(); this.loadEnv() }
      if (id === 'automations') this.loadAutomations()
      if (id === 'database') this.loadDBTabData()
      if (id === 'environment') {
        if (!this.environment) this.loadEnvironment()
        this.loadAPIKeys()
      }
    },
    get pageTitle() {
      const map = { dashboard:'总览', 'new-run':'新建任务', runs:'任务列表', automations:'定时任务', agents:'Agent 接入', database:'数据库', 'run-detail':'任务详情', environment:'环境检查', models:'模型管理' }
      return map[this.page] ?? ''
    },

    runSubjectSummary(run) {
      const subjects = run?.spec?.subjects ?? []
      if (subjects.length) return subjects.join(', ')
      const models = run?.spec?.models ?? []
      return models.length ? models.join(', ') : '—'
    },

    // ─── Agent 接入管理 ─────────────────────────────────────
    get dockerBackedFrameworkCount() {
      return (this.config?.frameworks ?? []).filter(f => f.sandbox_mode === 'docker' || f.sandbox_provider === 'docker').length
    },
    get filteredAgentSubjects() {
      const q = (this.agentSubjectQuery || '').trim().toLowerCase()
      const rows = this.config?.subjects ?? []
      if (!q) return rows
      return rows.filter(s => [s.id, s.framework, s.model, s.skill, s.kind].some(v => String(v || '').toLowerCase().includes(q)))
    },
    agentSubjectsByFramework(name) {
      return (this.config?.subjects ?? []).filter(s => s.framework === name)
    },
    agentSubjectsBySkill(name) {
      return (this.config?.subjects ?? []).filter(s => s.skill === name)
    },
    agentSkillsForFramework(name) {
      return (this.config?.skills ?? []).filter(skill => {
        const list = skill.compatible_frameworks || []
        return !list.length || list.includes(name)
      })
    },
    agentModelsLabel(framework) {
      const models = framework?.compatible_models || []
      if (!models.length) return '全部已启用模型'
      return models.join(', ')
    },
    agentLangsLabel(item) {
      const langs = item?.compatible_languages || []
      return langs.length ? langs.join(', ') : '全部语言'
    },
    agentKeyInfo(key) {
      const api = (this.apiKeys || []).find(k => k.key === key)
      if (api) return { known: true, set: !!api.value_set, label: api.label || key }
      const model = (this.models || []).find(m => m.api_key_env === key)
      if (model) return { known: true, set: !!model.api_key_set, label: model.name + ' Key' }
      return { known: false, set: false, label: key }
    },
    agentKeyStyle(key) {
      const info = this.agentKeyInfo(key)
      if (info.set) return 'background:var(--success-bg);color:var(--green)'
      if (info.known) return 'background:var(--error-bg);color:var(--red)'
      return 'background:var(--badge-bg);color:var(--fg-muted)'
    },
    agentKeyText(key) {
      const info = this.agentKeyInfo(key)
      if (info.set) return key + ' 已设置'
      if (info.known) return key + ' 未设置'
      return key + ' 按需透传'
    },
    agentFrameworkStatus(framework) {
      if (!framework) return '未知'
      if (this.config?.agents_config_error) return '配置异常'
      const keys = framework.env_from_host || []
      const missingKnown = keys.filter(key => {
        const info = this.agentKeyInfo(key)
        return info.known && !info.set
      })
      if (missingKnown.length) return '需补密钥'
      if ((framework.sandbox_mode || framework.sandbox_provider) === 'docker' && this.env && !this.env.docker_available) return '需 Docker'
      return '可运行'
    },
    agentFrameworkStatusStyle(framework) {
      const status = this.agentFrameworkStatus(framework)
      if (status === '可运行') return 'background:var(--success-bg);color:var(--green)'
      if (status === '需补密钥' || status === '需 Docker') return 'background:var(--warn-bg);color:var(--yellow)'
      return 'background:var(--error-bg);color:var(--red)'
    },
    copyText(text, label = '内容') {
      const done = () => this.showToast(label + '已复制', 'ok')
      if (navigator.clipboard?.writeText) {
        navigator.clipboard.writeText(text || '').then(done).catch(() => this.showToast('复制失败', 'err'))
        return
      }
      const input = document.createElement('textarea')
      input.value = text || ''
      document.body.appendChild(input)
      input.select()
      document.execCommand('copy')
      input.remove()
      done()
    },
    useAgentSubject(subject) {
      if (!subject) return
      this.form.combinations = [{
        _id: ++this._comboSeq,
        framework: subject.framework || 'model_api',
        model: subject.model || '',
        skill: subject.skill || 'no_skill',
      }]
      this.form.models = []
      if (subject.sandbox_mode === 'docker') this.form.use_docker = true
      this._comboKey++
      this.goto('new-run')
      this.showToast('已载入 ' + subject.id, 'ok')
    },
    useAgentFramework(framework) {
      const subject = this.agentSubjectsByFramework(framework?.name || '').find(s => s.kind === 'cli_agent') ||
        this.agentSubjectsByFramework(framework?.name || '')[0]
      if (subject) this.useAgentSubject(subject)
    },
    agentReadinessCounts() {
      const rows = this.config?.frameworks ?? []
      return {
        ready: rows.filter(f => this.agentFrameworkStatus(f) === '可运行').length,
        warning: rows.filter(f => ['需补密钥', '需 Docker'].includes(this.agentFrameworkStatus(f))).length,
        error: rows.filter(f => !['可运行', '需补密钥', '需 Docker'].includes(this.agentFrameworkStatus(f))).length,
      }
    },
    agentRuntimeCatalog() {
      const globs = 'generated_test.py\ntest_*.py\n*_test.go\n*Test.java\n*test*.cpp'
      return [
        {
          id: 'codex',
          title: 'Codex CLI',
          vendor: 'OpenAI',
          framework: 'codex',
          commandName: 'codex',
          badge: '需安装',
          summary: '把 Codex 作为独立 Agent Runtime 接入，安装进 Agent 沙箱镜像后再参与评测。',
          image: 'utbench-agent-base:latest',
          installCommand: 'npm install -g @openai/codex',
          installPackage: '@openai/codex',
          env_from_host: 'OPENAI_API_KEY',
          command: 'PROMPT="$(cat {{.ContainerPrompt}})" &&\ncodex exec --skip-git-repo-check --sandbox workspace-write \\\n  -c \'model_provider="utbench"\' \\\n  -c \'model_providers.utbench.name="{{.ModelProvider}} via UT-Bench"\' \\\n  -c \'model_providers.utbench.base_url="{{.ModelEndpoint}}"\' \\\n  -c \'model_providers.utbench.env_key="{{.ModelAPIKeyEnv}}"\' \\\n  -c \'model_providers.utbench.wire_api="chat"\' \\\n  --model "{{.ModelID}}" "$PROMPT"',
          output_globs: globs,
          needsInstall: true,
        },
        {
          id: 'kilo',
          title: 'Kilo Code',
          vendor: 'Kilo',
          framework: 'kilo',
          commandName: 'kilo',
          badge: '待适配',
          summary: '按产品运行时接入 Kilo，先把 CLI 安装进沙箱，再在高级适配里确认非交互命令。',
          image: 'utbench-agent-base:latest',
          installCommand: 'npm install -g @kilocode/cli',
          installPackage: '@kilocode/cli',
          env_from_host: 'KILO_API_KEY\nKILO_PROVIDER\nKILOCODE_MODEL\nKILO_ORG_ID',
          command: 'mkdir -p "$XDG_CONFIG_HOME/kilo" &&\nprintf \'{"permission":"allow"}\' > "$XDG_CONFIG_HOME/kilo/opencode.json" &&\nPROMPT="$(cat {{.ContainerPrompt}})" &&\nkilo run --auto --timeout 600 "$PROMPT"',
          output_globs: globs,
          needsInstall: true,
        },
        {
          id: 'claudecode',
          title: 'Claude Code',
          vendor: 'Anthropic',
          framework: 'claudecode',
          commandName: 'claude',
          badge: '镜像内置',
          summary: 'Agent 沙箱镜像已内置 Claude Code，配置密钥后即可组合模型与 Skill。',
          image: 'utbench-agent-base:latest',
          installCommand: 'npm install -g @anthropic-ai/claude-code',
          installPackage: '@anthropic-ai/claude-code',
          env_from_host: 'ANTHROPIC_API_KEY\nANTHROPIC_AUTH_TOKEN\nANTHROPIC_CUSTOM_HEADERS\nCLAUDE_CODE_USE_BEDROCK\nCLAUDE_CODE_USE_VERTEX\nCLAUDE_CODE_USE_FOUNDRY',
          command: 'mkdir -p "{{.ContainerWorkdir}}/.claude" &&\nPROMPT="$(cat {{.ContainerPrompt}})" &&\nclaude -p --output-format stream-json --verbose --permission-mode bypassPermissions --max-turns 50 --model "{{.ModelID}}" "$PROMPT"',
          output_globs: globs,
          bundled: true,
        },
        {
          id: 'opencode',
          title: 'OpenCode',
          vendor: 'SST',
          framework: 'opencode',
          commandName: 'opencode',
          badge: '镜像内置',
          summary: 'OpenCode 已作为通用 Agent Runtime 内置，适合 OpenAI-compatible 模型端点。',
          image: 'utbench-agent-base:latest',
          installCommand: 'npm install -g opencode-ai opencode-linux-x64-baseline',
          installPackage: 'opencode-ai opencode-linux-x64-baseline',
          env_from_host: 'DEEPSEEK_API_KEY\nDASHSCOPE_API_KEY\nMINIMAX_API_KEY\nARK_API_KEY\nBIGMODEL_API_KEY',
          command: 'PROMPT="$(cat {{.ContainerPrompt}})" &&\nopencode run --print-logs --dangerously-skip-permissions --model "{{.ModelID}}" "$PROMPT"',
          output_globs: globs,
          bundled: true,
        },
        {
          id: 'codebuddy',
          title: 'CodeBuddy',
          vendor: 'Tencent',
          framework: 'codebuddy',
          commandName: 'codebuddy',
          badge: '镜像内置',
          summary: 'CodeBuddy 已内置到 Agent 沙箱，适合用自定义模型端点做单测生成。',
          image: 'utbench-agent-base:latest',
          installCommand: 'npm install -g @tencent-ai/codebuddy-code',
          installPackage: '@tencent-ai/codebuddy-code',
          env_from_host: 'CODEBUDDY_API_KEY\nCODEBUDDY_INTERNET_ENVIRONMENT\nDEEPSEEK_API_KEY\nDASHSCOPE_API_KEY\nMINIMAX_API_KEY\nARK_API_KEY\nBIGMODEL_API_KEY',
          command: 'mkdir -p "{{.ContainerWorkdir}}/.codebuddy" &&\nPROMPT="$(cat {{.ContainerPrompt}})" &&\ncodebuddy -p -y --output-format stream-json --max-turns 50 "$PROMPT"',
          output_globs: globs,
          bundled: true,
        },
        {
          id: 'custom',
          title: '自定义 Agent Runtime',
          vendor: 'Custom',
          framework: 'custom-agent',
          commandName: 'agent-cli',
          badge: '高级',
          summary: '用于接入公司内部 Agent 或未预置产品，需要提供安装方式和非交互命令。',
          image: 'utbench-agent-base:latest',
          installCommand: '在 docker/agents/Dockerfile 安装你的 Agent CLI',
          installPackage: '',
          env_from_host: 'MY_AGENT_API_KEY',
          command: 'PROMPT="$(cat {{.ContainerPrompt}})" &&\nagent-cli run --model "{{.ModelID}}" --prompt "$PROMPT"',
          output_globs: globs,
          needsInstall: true,
          needsAdapter: true,
        },
      ]
    },
    agentRuntimeSpec(id = this.agentPreset) {
      return this.agentRuntimeCatalog().find(item => item.id === id) || this.agentRuntimeCatalog()[0]
    },
    agentFrameworkByName(name) {
      return (this.config?.frameworks ?? []).find(f => f.name === name)
    },
    agentRuntimeInstalled(id = this.agentPreset) {
      const spec = this.agentRuntimeSpec(id)
      return !!this.agentFrameworkByName(spec?.framework)
    },
    agentRuntimeStatusText(id = this.agentPreset) {
      const spec = this.agentRuntimeSpec(id)
      if (this.agentRuntimeInstalled(id)) return '已接入'
      if (spec?.needsAdapter) return '待适配'
      if (spec?.needsInstall) return '需安装'
      if (spec?.bundled) return '镜像内置'
      return '可接入'
    },
    agentRuntimeStatusStyle(id = this.agentPreset) {
      const text = this.agentRuntimeStatusText(id)
      if (text === '已接入' || text === '镜像内置') return 'background:var(--success-bg);color:var(--green)'
      if (text === '需安装' || text === '待适配') return 'background:var(--warn-bg);color:var(--yellow)'
      return 'background:var(--badge-bg);color:var(--fg-muted)'
    },
    agentRuntimeImagePresent() {
      return !!this.env?.agent_image_present
    },
    agentRuntimeImageText() {
      if (!this.env?.docker_available) return 'Docker 未就绪'
      return this.agentRuntimeImagePresent() ? 'Agent 镜像已构建' : 'Agent 镜像未构建'
    },
    agentRuntimeImageStyle() {
      if (this.agentRuntimeImagePresent()) return 'background:var(--success-bg);color:var(--green)'
      return 'background:var(--warn-bg);color:var(--yellow)'
    },
    agentCheckKey(id = this.agentPreset) {
      const spec = this.agentRuntimeSpec(id)
      return (spec?.image || 'utbench-agent-base:latest') + '::' + (spec?.commandName || '')
    },
    agentCheckResult(id = this.agentPreset) {
      return this.agentCheckResults[this.agentCheckKey(id)]
    },
    agentCheckText(id = this.agentPreset) {
      const result = this.agentCheckResult(id)
      if (!result) return '尚未检查'
      return result.ok ? 'CLI 已安装' : 'CLI 未安装'
    },
    agentCheckStyle(id = this.agentPreset) {
      const result = this.agentCheckResult(id)
      if (!result) return 'background:var(--badge-bg);color:var(--fg-muted)'
      return result.ok ? 'background:var(--success-bg);color:var(--green)' : 'background:var(--error-bg);color:var(--red)'
    },
    agentCheckMessage(id = this.agentPreset) {
      const result = this.agentCheckResult(id)
      if (!result) return ''
      if (result.ok) return result.version || result.output || 'agent cli found'
      return result.error || result.output || 'agent cli not found'
    },
    agentCanInstallCLI() {
      const spec = this.agentRuntimeSpec()
      return !!spec?.installPackage && !this.agentInstallingCLI && this.env?.docker_available
    },
    agentInstallCLIText() {
      if (this.agentInstallingCLI) return '安装中...'
      const spec = this.agentRuntimeSpec()
      if (!spec?.installPackage) return '未配置安装包'
      return '安装 CLI 到镜像'
    },
    async installAgentCLI() {
      const spec = this.agentRuntimeSpec()
      if (!spec?.installPackage) { this.showToast('这个 Agent 还没有配置可自动安装的 CLI 包', 'warn'); return }
      if (!this.env?.docker_available) { this.showToast('Docker 还不可用，无法安装 CLI', 'warn'); return }
      this.agentInstallingCLI = true
      this.agentFormError = ''
      try {
        const r = await fetch('/api/agents/install-cli', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ runtime: spec.id, image: this.agentForm.image || spec.image, package: spec.installPackage }),
        })
        const data = await r.json()
        if (!r.ok) throw new Error(data.error || 'HTTP ' + r.status)
        this.buildModalOpen = true
        this.buildTarget = data.target || 'agent'
        this.buildId = data.build_id || ''
        this.buildStatus = data.status || 'pending'
        this.buildError = data.error || ''
        this.buildLogs = []
        if (this.buildId) this.startBuildSSE(this.buildId)
        this.showToast('开始安装 ' + spec.title + ' CLI 到 Agent 镜像', 'ok')
      } catch (e) {
        this.agentFormError = e.message || String(e)
        this.showToast('安装 Agent CLI 失败：' + this.agentFormError, 'err', 6000)
      } finally {
        this.agentInstallingCLI = false
      }
    },
    async checkAgentRuntime() {
      const spec = this.agentRuntimeSpec()
      if (!spec?.commandName) { this.showToast('没有可检查的 Agent CLI', 'warn'); return }
      if (!this.agentRuntimeImagePresent()) { this.showToast('Agent 镜像还没有构建', 'warn'); return }
      this.agentChecking = true
      this.agentFormError = ''
      try {
        const r = await fetch('/api/agents/check', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ image: this.agentForm.image || spec.image, command: spec.commandName }),
        })
        const data = await r.json()
        if (!r.ok) throw new Error(data.error || 'HTTP ' + r.status)
        this.agentCheckResults = { ...this.agentCheckResults, [this.agentCheckKey()]: data }
        if (data.ok) this.showToast(spec.title + ' CLI 已安装：' + (data.version || data.command), 'ok')
        else this.showToast(spec.title + ' CLI 未安装：' + (data.error || data.output || 'not found'), 'err', 6000)
      } catch (e) {
        this.agentFormError = e.message || String(e)
        this.showToast('检查 Agent 失败：' + this.agentFormError, 'err', 6000)
      } finally {
        this.agentChecking = false
      }
    },
    selectAgentRuntime(id) {
      this.agentFormError = ''
      this.agentRuntimeInstallConfirmed = false
      this.applyAgentPreset(id)
    },
    openAgentWizard() {
      this.agentFormError = ''
      this.agentWizardOpen = true
      this.selectAgentRuntime(this.agentPreset || 'codex')
    },
    closeAgentWizard() { this.agentWizardOpen = false },
    applyAgentPreset(preset) {
      this.agentPreset = preset
      const spec = this.agentRuntimeSpec(preset)
      const next = {
        name: spec.framework,
        image: spec.image || 'utbench-agent-base:latest',
        timeout_seconds: 600,
        network_disabled: false,
        compatible_languages: ['python', 'go', 'java', 'cpp'],
        compatible_models: [],
        output_globs: spec.output_globs || 'generated_test.py\ntest_*.py\n*_test.go\n*Test.java\n*test*.cpp',
        env_from_host: spec.env_from_host || '',
        command: spec.command || '',
      }
      this.agentAdvancedOpen = !!spec.needsAdapter
      this.agentForm = { ...this.agentForm, ...next }
    },
    agentCanInstallSelected() {
      const spec = this.agentRuntimeSpec()
      if (!spec || this.agentSaving || this.agentRuntimeInstalled(spec.id)) return false
      if (!String(this.agentForm.name || '').trim()) return false
      if (!String(this.agentForm.command || '').trim()) return false
      return true
    },
    toggleAgentLanguage(lang) {
      const set = new Set(this.agentForm.compatible_languages || [])
      if (set.has(lang)) set.delete(lang)
      else set.add(lang)
      this.agentForm.compatible_languages = ['python', 'go', 'java', 'cpp'].filter(x => set.has(x))
    },
    toggleAgentModel(name) {
      const set = new Set(this.agentForm.compatible_models || [])
      if (set.has(name)) set.delete(name)
      else set.add(name)
      this.agentForm.compatible_models = Array.from(set)
    },
    splitLines(value) {
      return String(value || '').split(/\r?\n|,/).map(s => s.trim()).filter(Boolean)
    },
    async saveAgentFramework() {
      this.agentFormError = ''
      const spec = this.agentRuntimeSpec()
      if (this.agentRuntimeInstalled(spec?.id)) {
        this.showToast((spec?.title || 'Agent') + ' 已经接入，无需重复安装', 'warn')
        this.agentWizardOpen = false
        return
      }
      const f = this.agentForm
      if (!f.name.trim()) { this.agentFormError = '请填写 Agent 名称'; return }
      if (!f.command.trim()) { this.agentFormError = '请填写启动命令'; return }
      if ((spec?.needsInstall || spec?.needsAdapter) && !this.agentRuntimeInstallConfirmed) {
        this.agentFormError = '请先确认该 Agent CLI 已安装到沙箱镜像，并且高级适配命令可非交互运行'
        return
      }
      this.agentSaving = true
      try {
        const payload = {
          name: f.name,
          command: f.command,
          image: f.image,
          timeout_seconds: Number(f.timeout_seconds) || 600,
          network_disabled: !!f.network_disabled,
          env_from_host: this.splitLines(f.env_from_host),
          compatible_models: f.compatible_models || [],
          compatible_languages: f.compatible_languages || [],
          output_globs: this.splitLines(f.output_globs),
        }
        const r = await fetch('/api/agents/frameworks', { method:'POST', headers:{'Content-Type':'application/json'}, body:JSON.stringify(payload) })
        const data = await r.json()
        if (!r.ok) throw new Error(data.error || 'HTTP ' + r.status)
        this.showToast('Agent Runtime 已接入：' + data.name, 'ok')
        this.agentWizardOpen = false
        await this.loadConfig()
      } catch (e) {
        this.agentFormError = e.message || String(e)
        this.showToast('接入 Agent 失败：' + this.agentFormError, 'err', 6000)
      } finally {
        this.agentSaving = false
      }
    },

    openSkillModal() {
      this.skillFormError = ''
      this.skillPackages = []
      this.skillInstallMode = 'scan'
      this.skillCommand = ''
      this.skillCommandOutput = ''
      this.skillUploadFileName = ''
      this.skillForm = { skill_root: '', names: [] }
      this.skillModalOpen = true
    },

    closeSkillModal() { this.skillModalOpen = false },

    toggleSkillPackage(name) {
      const idx = this.skillForm.names.indexOf(name)
      if (idx >= 0) this.skillForm.names.splice(idx, 1)
      else this.skillForm.names.push(name)
    },

    async scanSkillPackages() {
      this.skillScanning = true
      this.skillFormError = ''
      try {
        const r = await fetch('/api/agents/skills/scan?' + new URLSearchParams({ skill_root: this.skillForm.skill_root || './skills' }), { cache: 'no-store' })
        const data = await r.json()
        if (!r.ok) throw new Error(data.error || 'HTTP ' + r.status)
        this.skillPackages = data.packages || []
        if (!this.skillPackages.length) this.skillFormError = '目录下未发现 Skill 包'
      } catch (e) {
        this.skillFormError = e.message || String(e)
        this.showToast('扫描失败：' + this.skillFormError, 'err')
      } finally {
        this.skillScanning = false
      }
    },

    async uploadSkillPackage(evt) {
      const file = evt?.target?.files?.[0]
      if (!file) return
      this.skillUploading = true
      this.skillFormError = ''
      this.skillUploadFileName = file.name
      try {
        const body = new FormData()
        body.append('skill_root', this.skillForm.skill_root || './skills')
        body.append('file', file)
        const r = await fetch('/api/agents/skills/upload', { method: 'POST', body })
        const data = await r.json()
        if (!r.ok) throw new Error(data.error || 'HTTP ' + r.status)
        this.skillPackages = data.packages || []
        this.skillForm.skill_root = data.skill_root || this.skillForm.skill_root || './skills'
        this.showToast('已上传并解压 Skill 包：' + file.name, 'ok')
      } catch (e) {
        this.skillFormError = e.message || String(e)
        this.showToast('上传失败：' + this.skillFormError, 'err', 6000)
      } finally {
        this.skillUploading = false
        if (evt?.target) evt.target.value = ''
      }
    },

    async runSkillInstallCommand() {
      if (!this.skillCommand.trim()) return
      this.skillCommandRunning = true
      this.skillFormError = ''
      this.skillCommandOutput = '$ ' + this.skillCommand + '\n'
      try {
        const r = await fetch('/api/agents/skills/command', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ skill_root: this.skillForm.skill_root || './skills', command: this.skillCommand }),
        })
        const data = await r.json()
        if (!r.ok) throw new Error(data.error || 'HTTP ' + r.status)
        this.skillCommandOutput += data.output || '(命令执行完成，无输出)'
        this.skillPackages = data.packages || []
        this.skillForm.skill_root = data.skill_root || this.skillForm.skill_root || './skills'
        this.showToast('命令执行完成，已刷新 Skill 包列表', 'ok')
      } catch (e) {
        this.skillFormError = e.message || String(e)
        this.skillCommandOutput += '\n' + this.skillFormError
        this.showToast('命令执行失败：' + this.skillFormError, 'err', 6000)
      } finally {
        this.skillCommandRunning = false
      }
    },

    async installSkills() {
      if (!this.skillForm.names.length) return
      this.skillSaving = true
      this.skillFormError = ''
      try {
        const r = await fetch('/api/agents/skills', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ skill_root: this.skillForm.skill_root || './skills', names: this.skillForm.names }),
        })
        const data = await r.json()
        if (!r.ok) throw new Error(data.error || 'HTTP ' + r.status)
        this.showToast('已安装 ' + (data.installed || []).length + ' 个 Skill', 'ok')
        this.skillModalOpen = false
        await this.loadConfig()
      } catch (e) {
        this.skillFormError = e.message || String(e)
        this.showToast('安装失败：' + this.skillFormError, 'err', 6000)
      } finally {
        this.skillSaving = false
      }
    },

    toggleSkillSelection(name) {
      const idx = this.selectedSkillNames.indexOf(name)
      if (idx >= 0) this.selectedSkillNames.splice(idx, 1)
      else this.selectedSkillNames.push(name)
    },

    viewSkillDetail(skill) {
      this.skillDetail = skill
      this.skillDetailOpen = true
    },

    compareSelectedSkills() {
      if (this.selectedSkillNames.length < 2) return
      const skills = (this.config?.skills ?? []).filter(s => this.selectedSkillNames.includes(s.name))
      this.compareSkills = skills
      this.skillCompareOpen = true
    },

    get skillCompareFields() {
      return [
        { key: 'version',     label: '版本',              value: s => s.version || '-' },
        { key: 'inject_mode', label: '注入模式',           value: s => s.inject_mode || '-' },
        { key: 'description', label: '描述',              value: s => s.description || '-' },
        { key: 'instruction_path', label: '指令路径',    value: s => s.instruction_path || '-' },
        { key: 'files',       label: '关联文件',           value: s => (s.files || []) },
        { key: 'compatible_frameworks', label: '兼容 Framework', value: s => (s.compatible_frameworks || []) },
        { key: 'compatible_languages',  label: '兼容语言', value: s => (s.compatible_languages || []) },
      ]
    },

    // ─── 模型管理 ─────────────────────────────────────────────
    async loadModels() {
      this.modelsLoading = true
      try {
        const r = await fetch('/api/models', { cache: 'no-store' })
        if (!r.ok) throw new Error('HTTP ' + r.status)
        this.models = await r.json()
      } catch (e) { this.showToast('加载模型失败：' + e.message, 'err') }
      finally { this.modelsLoading = false }
    },

    openAddModel() {
      this.modelFormMode = 'create'
      this.modelForm = {
        name: '', enabled: true, provider: '', model_id: '',
        api_endpoint: '', anthropic_endpoint: '',
        api_key_env: '', api_key: '', api_key_set: false,
        parameters: { temperature: 0.7, top_p: 0.9, max_tokens: 4096 },
      }
      this.modelFormError = ''
      this.modelFormOpen = true
    },

    openEditModel(m) {
      this.modelFormMode = 'edit'
      this.modelForm = {
        name: m.name, enabled: !!m.enabled, provider: m.provider || '',
        model_id: m.model_id || '', api_endpoint: m.api_endpoint || '',
        anthropic_endpoint: m.anthropic_endpoint || '',
        api_key_env: m.api_key_env || '', api_key: '', api_key_set: !!m.api_key_set,
        parameters: Object.assign({ temperature: 0.7, top_p: 0.9, max_tokens: 4096 }, m.parameters || {}),
      }
      this.modelFormError = ''
      this.modelFormOpen = true
    },

    closeModelForm() { this.modelFormOpen = false },

    async saveModel() {
      this.modelFormError = ''
      const f = this.modelForm
      if (!f.name.trim()) { this.modelFormError = '模型名称必填'; return }
      if (!f.provider.trim()) { this.modelFormError = '提供商必填'; return }
      if (!f.model_id.trim()) { this.modelFormError = '模型 ID 必填'; return }
      if (!f.api_endpoint.trim()) { this.modelFormError = 'API 端点必填'; return }
      const body = { ...f }
      try {
        let r
        if (this.modelFormMode === 'create') {
          r = await fetch('/api/models', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body) })
        } else {
          r = await fetch('/api/models/' + encodeURIComponent(f.name), { method: 'PUT', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body) })
        }
        const data = await r.json()
        if (!r.ok) { this.modelFormError = data.error || '保存失败'; return }
        this.modelFormOpen = false
        this.showToast(this.modelFormMode === 'create' ? '模型已添加' : '模型已更新', 'ok')
        await this.loadModels()
        await this.loadConfig() // 刷新新建任务页面的可选模型列表
      } catch (e) { this.modelFormError = String(e) }
    },

    async deleteModel(name) {
      if (!confirm(`确定删除模型 "${name}"？此操作将从 models.yaml 移除该条目。`)) return
      try {
        const r = await fetch('/api/models/' + encodeURIComponent(name), { method: 'DELETE' })
        if (!r.ok) { const d = await r.json(); throw new Error(d.error || 'HTTP ' + r.status) }
        this.models = this.models.filter(m => m.name !== name)
        this.form.models = this.form.models.filter(m => m !== name)
        const nextResults = { ...this.modelTestResults }
        delete nextResults[name]
        this.modelTestResults = nextResults
        this.showToast('模型 "' + name + '" 已删除', 'ok')
        await this.loadModels()
        await this.loadConfig()
      } catch (e) { this.showToast('删除失败：' + e.message, 'err') }
    },

    isModelTesting(name) {
      return !!this.modelTesting[name]
    },

    modelTestResult(name) {
      return this.modelTestResults[name] || null
    },

    async testModel(name) {
      this.modelTesting = { ...this.modelTesting, [name]: true }
      try {
        const r = await fetch('/api/models/' + encodeURIComponent(name) + '/test', { method: 'POST' })
        const data = await r.json()
        if (!r.ok) throw new Error(data.error || 'HTTP ' + r.status)
        this.modelTestResults = { ...this.modelTestResults, [name]: data }
        this.showToast(data.ok ? `模型 "${name}" 连接正常` : `模型 "${name}" 连接失败：${data.message}`, data.ok ? 'ok' : 'err', 5000)
      } catch (e) {
        const data = { name, ok: false, status: 'request_error', message: e.message || String(e), checked_at: new Date().toISOString() }
        this.modelTestResults = { ...this.modelTestResults, [name]: data }
        this.showToast(`模型 "${name}" 连接失败：${data.message}`, 'err', 5000)
      } finally {
        const next = { ...this.modelTesting }
        delete next[name]
        this.modelTesting = next
      }
    },

    async testAllModels() {
      if (!this.models.length) return
      this.modelTestingAll = true
      const testing = {}
      for (const m of this.models) testing[m.name] = true
      this.modelTesting = { ...this.modelTesting, ...testing }
      try {
        const r = await fetch('/api/models/test-all', { method: 'POST' })
        const data = await r.json()
        if (!r.ok) throw new Error(data.error || 'HTTP ' + r.status)
        const next = { ...this.modelTestResults }
        for (const item of data) next[item.name] = item
        this.modelTestResults = next
        const ok = data.filter(item => item.ok).length
        this.showToast(`连接测试完成：${ok}/${data.length} 个正常`, ok === data.length ? 'ok' : 'warn', 6000)
      } catch (e) {
        this.showToast('一键测试失败：' + (e.message || String(e)), 'err', 6000)
      } finally {
        this.modelTestingAll = false
        this.modelTesting = {}
      }
    },

    async rerunRun(runID) {
      if (!confirm('重新提交该任务？会以相同配置创建一个新的 run。')) return
      try {
        const r = await fetch('/api/runs/' + encodeURIComponent(runID) + '/rerun', { method: 'POST' })
        const data = await r.json()
        if (!r.ok) { this.showToast('重跑失败：' + (data.error || 'HTTP ' + r.status), 'err'); return }
        const newRunId = data.run_id || data.runId || data.id
        if (!newRunId) { this.showToast('重跑失败：backend did not return run_id', 'err'); return }
        this.showToast('已创建 Docker 重跑任务 ' + newRunId, 'ok')
        await this.openRun(newRunId)
        this.loadRuns()
      } catch (e) { this.showToast('重跑失败：' + e.message, 'err') }
    },

    async deleteRun(runID) {
      if (!confirm('确定删除该任务？会删除磁盘上所有产物（生成结果、评测、报告），不可恢复。')) return
      try {
        const r = await fetch('/api/runs/' + encodeURIComponent(runID), { method: 'DELETE' })
        const data = await r.json()
        if (!r.ok) { this.showToast('删除失败：' + (data.error || 'HTTP ' + r.status), 'err'); return }
        this.showToast('已删除 ' + runID, 'ok')
        await this.loadRuns()
        if (this.currentRun && this.currentRun.run_id === runID) {
          this.currentRun = null; this.goto('runs')
        }
      } catch (e) { this.showToast('删除失败：' + e.message, 'err') }
    },

    async renameRun(runID, currentLabel) {
      const label = prompt('输入任务备注名（留空清除）：', currentLabel || '')
      if (label === null) return
      try {
        const r = await fetch('/api/runs/' + encodeURIComponent(runID), {
          method: 'PATCH',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ label: label.trim() })
        })
        const data = await r.json()
        if (!r.ok) { this.showToast('重命名失败：' + (data.error || 'HTTP ' + r.status), 'err'); return }
        this.showToast('已更新备注名', 'ok')
        await this.loadRuns()
        await this.loadAssets()
        if (this.currentRun && this.currentRun.run_id === runID) {
          this.currentRun.label = label.trim()
        }
      } catch (e) { this.showToast('重命名失败：' + e.message, 'err') }
    },

    async reevaluateRun(runID) {
      if (!confirm('重新运行 evaluator？会重新编译、运行测试、覆盖率和变异测试，并覆盖当前 evaluation_result.json。')) return
      this.runActionBusy = 'reevaluate'
      try {
        const r = await fetch('/api/runs/' + encodeURIComponent(runID) + '/reevaluate', { method: 'POST' })
        const data = await r.json()
        if (!r.ok) { this.showToast('重新评测失败：' + (data.error || 'HTTP ' + r.status), 'err', 6000); return }
        this.currentReport = null
        this.showToast('evaluator 重新运行完成：' + data.evaluation_path, 'ok', 6000)
        await this.refreshDetail()
      } catch (e) {
        this.showToast('重新评测失败：' + e.message, 'err', 6000)
      } finally {
        this.runActionBusy = ''
      }
    },

    async regenerateReport(runID) {
      const defaultPath = `artifacts/runs/${runID}/evaluation/evaluation_result.json`
      const evaluationPath = prompt('输入 evaluation_result.json 路径：', defaultPath)
      if (evaluationPath === null) return
      this.runActionBusy = 'report'
      try {
        const r = await fetch('/api/runs/' + encodeURIComponent(runID) + '/regenerate-report', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ evaluation_path: evaluationPath.trim() || defaultPath }),
        })
        const data = await r.json()
        if (!r.ok) { this.showToast('生成报告失败：' + (data.error || 'HTTP ' + r.status), 'err', 6000); return }
        this.showToast('报告已重新生成', 'ok')
        this.detailTab = 'report'
        await this.loadReport()
      } catch (e) {
        this.showToast('生成报告失败：' + e.message, 'err', 6000)
      } finally {
        this.runActionBusy = ''
      }
    },

    async runControl(runID, action) {
      const actionLabels = { pause:'暂停', resume:'继续', cancel:'取消' }
      try {
        const r = await fetch('/api/runs/' + encodeURIComponent(runID) + '/' + action, { method: 'POST' })
        const data = await r.json()
        if (!r.ok) { this.showToast((actionLabels[action] || action) + '失败：' + (data.error || 'HTTP ' + r.status), 'err'); return }
        this.showToast((actionLabels[action] || action) + '成功', 'ok')
        await this.loadRuns()
        if (this.page === 'run-detail' && this.currentRun && this.currentRun.run_id === runID) await this.refreshDetail()
      } catch (e) { this.showToast((actionLabels[action] || action) + '失败：' + e.message, 'err') }
    },

    async toggleModelEnabled(m) {
      try {
        const body = { ...m, enabled: !m.enabled, api_key: '' }
        const r = await fetch('/api/models/' + encodeURIComponent(m.name), { method: 'PUT', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body) })
        if (!r.ok) { const d = await r.json(); throw new Error(d.error || 'HTTP ' + r.status) }
        await this.loadModels()
      } catch (e) { this.showToast('切换失败：' + e.message, 'err') }
    },

    applyTheme() {
      document.documentElement.dataset.theme = this.theme
      localStorage.setItem('utbench-theme', this.theme)
    },
    toggleTheme() {
      this.theme = this.theme === 'dark' ? 'light' : 'dark'
      this.applyTheme()
    },

    _escapeHtml(s) { return String(s).replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;') },
    _modelRows() {
      return (this.models && this.models.length) ? this.models : (this.config?.models ?? [])
    },
    _comboFrameworks() {
      const seen = new Set()
      const fws = []
      // model_api 作为 baseline 始终在第一位
      fws.push({ value: 'model_api', label: 'model_api（纯 API）' }); seen.add('model_api')
      for (const fw of (this.config?.frameworks ?? [])) {
        if (!seen.has(fw.name)) { fws.push({ value: fw.name, label: fw.name }); seen.add(fw.name) }
      }
      return fws
    },
    _comboModels(framework) {
      const all = this._modelRows().map(m => m.name)
      const fw = (this.config?.frameworks ?? []).find(f => f.name === framework)
      if (!fw || !fw.compatible_models || fw.compatible_models.length === 0) return all
      const allowed = new Set(fw.compatible_models)
      const filtered = all.filter(m => allowed.has(m))
      return filtered.length > 0 ? filtered : all
    },
    _comboSkills(framework) {
      const seen = new Set(['no_skill'])
      const skills = ['no_skill']
      for (const sk of (this.config?.skills ?? [])) {
        if (seen.has(sk.name)) continue
        if (!sk.compatible_frameworks || sk.compatible_frameworks.length === 0) {
          skills.push(sk.name); seen.add(sk.name)
        } else if (sk.compatible_frameworks.includes(framework)) {
          skills.push(sk.name); seen.add(sk.name)
        }
      }
      return skills
    },
    comboModelOptionLabel(name) {
      return String(name || '')
    },
    renderFrameworkOptions(current) {
      return this._comboFrameworks().map(fw =>
        `<option value="${this._escapeHtml(fw.value)}" ${fw.value===current?'selected':''}>${this._escapeHtml(fw.label)}</option>`
      ).join('')
    },
    renderModelOptions(framework, current) {
      const models = this._comboModels(framework)
      if (!models.includes(current) && models.length) current = models[0]
      return models.map(m =>
        `<option value="${this._escapeHtml(m)}" ${m===current?'selected':''}>${this._escapeHtml(this.comboModelOptionLabel(m))}</option>`
      ).join('')
    },
    renderSkillOptions(framework, current) {
      const skills = this._comboSkills(framework)
      if (!skills.includes(current)) current = 'no_skill'
      return skills.map(s =>
        `<option value="${this._escapeHtml(s)}" ${s===current?'selected':''}>${this._escapeHtml(s)}</option>`
      ).join('')
    },
    renderCombinationRows() {
      const rows = this.form.combinations || []
      return rows.map((combo, idx) => {
        const id = combo._id ?? idx
        const disabled = rows.length <= 1 ? 'disabled' : ''
        return `<div class="flex items-center gap-2" data-combo-row="${this._escapeHtml(id)}">
          <select class="input-base text-[12px] flex-1" data-combo-id="${this._escapeHtml(id)}" data-combo-idx="${idx}" data-combo-type="framework">${this.renderFrameworkOptions(combo.framework)}</select>
          <select class="input-base text-[12px] flex-1" data-combo-id="${this._escapeHtml(id)}" data-combo-idx="${idx}" data-combo-type="model">${this.renderModelOptions(combo.framework, combo.model)}</select>
          <select class="input-base text-[12px] flex-1" data-combo-id="${this._escapeHtml(id)}" data-combo-idx="${idx}" data-combo-type="skill">${this.renderSkillOptions(combo.framework, combo.skill)}</select>
          <span class="font-mono text-[11px] px-2 py-1 rounded" style="background:var(--bg-overlay);color:var(--fg-subtle)">${this._escapeHtml(this.buildSubjectId(combo))}</span>
          <button type="button" data-combo-remove="${this._escapeHtml(id)}" class="text-[12px] px-1.5 py-0.5 rounded" style="color:var(--fg-muted);background:var(--bg-overlay)" ${disabled}>&times;</button>
        </div>`
      }).join('')
    },
    handleComboRowsChange(event) {
      const sel = event.target?.closest?.('select[data-combo-id]')
      if (!sel) return
      this.setCombinationValue(sel.dataset.comboId, Number(sel.dataset.comboIdx || 0), sel.dataset.comboType, sel.value)
    },
    handleComboRowsClick(event) {
      const btn = event.target?.closest?.('button[data-combo-remove]')
      if (!btn || btn.disabled) return
      const id = btn.dataset.comboRemove
      const idx = (this.form.combinations || []).findIndex(c => String(c._id) === String(id))
      if (idx >= 0 && this.form.combinations.length > 1) this.form.combinations.splice(idx, 1)
    },
    buildSubjectId(combo) {
      if (!combo.model) return ''
      // 与后端 agentconfig.sanitize 保持一致：保留 . 和 -，trim 下划线
      const sanitize = s => {
        let out = s.toLowerCase().trim().replace(/[^a-z0-9_.-]/g, '_').replace(/^_+|_+$/g, '')
        return out || 'unknown'
      }
      const fw = sanitize(combo.framework || 'model_api')
      const model = sanitize(combo.model)
      const skill = sanitize(combo.skill || 'no_skill')
      return `${fw}__${model}__${skill}`
    },
    addCombination() {
      if (this._comboAddLocked) return
      this._comboAddLocked = true
      setTimeout(() => { this._comboAddLocked = false }, 0)
      const models = this._comboModels('model_api')
      const defaultModel = models.includes('deepseek-v4-flash') ? 'deepseek-v4-flash' : (models[0] || '')
      this.form.models = []
      this.form.combinations.push({ _id: ++this._comboSeq, framework: 'model_api', model: defaultModel, skill: 'no_skill' })
    },
    setCombinationValue(rowID, idx, type, value) {
      let combo = null
      if (rowID !== undefined && rowID !== null && rowID !== '') {
        combo = (this.form.combinations || []).find(c => String(c._id) === String(rowID))
      }
      combo = combo || this.form.combinations[idx]
      if (!combo) return
      this.form.models = []
      if (type === 'framework') {
        combo.framework = value
        this.onCombinationFrameworkChange(this.form.combinations.indexOf(combo))
      } else if (type === 'model') {
        combo.model = value
      } else if (type === 'skill') {
        combo.skill = value
      }
    },
    onCombinationFrameworkChange(idx) {
      const combo = this.form.combinations[idx]
      const models = this._comboModels(combo.framework)
      if (!models.includes(combo.model)) {
        combo.model = models[0] || ''
      }
      const skills = this._comboSkills(combo.framework)
      if (!skills.includes(combo.skill)) {
        combo.skill = 'no_skill'
      }
    },
    syncPageVisibility() {
      const pages = ['dashboard', 'new-run', 'runs', 'automations', 'agents', 'database', 'environment', 'models', 'run-detail']
      const modals = ['model-form-modal', 'agent-form-modal', 'skill-form-modal', 'build-image-modal']
      requestAnimationFrame(() => {
        for (const id of pages) {
          const node = document.getElementById('partial-' + id)
          if (!node) continue
          const active = this.page === id
          node.style.display = active ? '' : 'none'
          if (active && node.firstElementChild) {
            node.firstElementChild.style.display = ''
          }
        }
        for (const id of modals) {
          const node = document.getElementById('partial-' + id)
          if (!node) continue
          node.style.display = ''
        }
      })
    },

    makeClientRunID() {
      const d = new Date()
      const pad = (n, w = 2) => String(n).padStart(w, '0')
      const nanos = pad(d.getUTCMilliseconds(), 3) + pad(Math.floor(Math.random() * 1000000), 6)
      return `${d.getUTCFullYear()}${pad(d.getUTCMonth() + 1)}${pad(d.getUTCDate())}T${pad(d.getUTCHours())}${pad(d.getUTCMinutes())}${pad(d.getUTCSeconds())}.${nanos}Z`
    },

    optimisticRunFromPayload(payload, runId) {
      return {
        run_id: runId,
        status: 'pending',
        started_at: new Date().toISOString(),
        use_docker: !!payload.use_docker,
        spec: {
          run_id: runId,
          models: payload.models || [],
          subjects: payload.subjects || [],
          languages: payload.languages || [],
          dataset_classes: payload.class ? String(payload.class).split(',').map(s => s.trim()).filter(Boolean) : [],
          dataset_scenario: payload.scenario || '',
          dataset_level: payload.level || '',
          max_samples: payload.max_samples,
          workers: payload.workers,
          reuse_generated: !!payload.reuse_generated,
          reuse_evaluation: !!payload.reuse_evaluation,
          mutation_enabled: !!payload.mutation_enabled,
        },
      }
    },

    async submitRun() {
      // 防双击：已经在提交中 / 已经跳到详情页就不要再发
      if (this.formSubmitting) return
      this.formError = ''
      // 将 combinations 转换为 subjects 数组
      this.form.subjects = (this.form.combinations || [])
        .filter(c => c.model)
        .map(c => this.buildSubjectId(c))
      // 只有 generate 和 full 阶段需要被测对象和语言选择
      if (this.form.phase === 'generate' || this.form.phase === 'full') {
        if (!this.selectedExecutionTargetCount) { this.formError = '请至少选择一个 subject 或模型'; return }
        if (!this.form.languages.length) { this.formError = '请至少选择一种语言'; return }
      }
      // evaluate 阶段需要数据源（source_run_id 或 manifest_path）
      if (this.form.phase === 'evaluate') {
        if (!this.form.source_run_id && !this.form.manifest_path) {
          this.formError = '请选择来源任务或填写 manifest 路径'; return
        }
      }
      // report 阶段需要数据源（source_run_id 或 evaluation_path）
      if (this.form.phase === 'report') {
        if (!this.form.source_run_id && !this.form.evaluation_path) {
          this.formError = '请选择来源任务或填写评测结果路径'; return
        }
      }
      this.formSubmitting = true
      const payload = JSON.parse(JSON.stringify(this.form))
      const optimisticRunId = String(payload.run_id || '').trim() || this.makeClientRunID()
      payload.run_id = optimisticRunId
      this.stopSSE()
      this.currentReport = null
      this.currentLogs = []
      this._logCount = 0
      this.detailTab = 'logs'
      this._resetLogPre()
      this.currentRun = this.optimisticRunFromPayload(payload, optimisticRunId)
      this.page = 'run-detail'
      this.syncPageVisibility()
      try {
        const r = await fetch('/api/runs', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(payload) })
        const data = await r.json()
        if (!r.ok) {
          const message = data.error || '提交失败'
          this.formError = message
          if (this.currentRun?.run_id === optimisticRunId) {
            this.currentRun.status = 'failed'
            this.currentRun.error = message
          }
          this.showToast('任务提交失败：' + message, 'err', 6000)
          return
        }
        const runId = data.run_id || data.runId || data.id || optimisticRunId
        if (!runId) { this.formError = 'backend did not return run_id'; return }
        try {
          await this.openRun(runId)
        } catch (err) {
          console.error('[submitRun] openRun failed:', err)
        }
        this.loadRuns().then(() => {
          if (this.page === 'run-detail' && this.currentRun?.run_id === runId) {
            const found = this.runs.find(r => r.run_id === runId)
            if (found) this.currentRun = { ...found }
          }
        }).catch(err => console.error('[submitRun] loadRuns failed:', err))
      } catch(e) {
        console.error('[submitRun] failed:', e)
        this.formError = String(e)
        if (this.currentRun?.run_id === optimisticRunId) {
          this.currentRun.status = 'failed'
          this.currentRun.error = String(e)
        }
        this.showToast('任务提交失败：' + String(e), 'err', 6000)
      }
      finally { this.formSubmitting = false }
    },

    async openRun(runId) {
      this.stopSSE(); this.currentReport = null; this.currentLogs = []; this._logCount = 0; this.detailTab = 'logs'
      this._resetLogPre()
      this.page = 'run-detail'
      this.syncPageVisibility()
      const found = this.runs.find(r => r.run_id === runId)
      this.currentRun = found ? { ...found } : { run_id: runId, status: 'pending', started_at: new Date().toISOString() }
      this.startSSE(runId)
    },

    startSSE(runId) {
      this.sseSource = new EventSource(`/api/runs/${runId}/events`)
      this.sseSource.onmessage = (e) => {
        try {
          const obj = JSON.parse(e.data)
          if (obj.type === 'snapshot') {
            // 后端在订阅瞬间把已缓冲的所有日志一次性发回，单事件 → 单 DOM 写入
            const lines = Array.isArray(obj.payload) ? obj.payload : []
            this.currentLogs = lines.slice(-this.LOG_DOM_CAP)
            this._applyLogSnapshotToDOM(this.currentLogs)
          }
          else if (obj.type === 'log') {
            this._enqueueLog(obj.payload)
          }
          else if (obj.type === 'done') {
            this.stopSSE()
            if (this.currentRun) this.currentRun.status = obj.payload
            this._flushLogs()
            this.loadRuns(); this.loadDatabase()
          }
        } catch {}
      }
      this.sseSource.onerror = () => this.stopSSE()
    },

    _applyLogSnapshotToDOM(lines) {
      // 等待 pre 元素挂载（首次打开 run-detail 时 x-if 还没生效）
      const tryApply = () => {
        const pre = document.getElementById('log-pre')
        if (!pre) { requestAnimationFrame(tryApply); return }
        pre.textContent = (lines || []).join('\n') + (lines && lines.length ? '\n' : '')
        const area = document.getElementById('log-area')
        if (area) area.scrollTop = area.scrollHeight
      }
      tryApply()
    },

    // ---- 日志高性能管线：避开 Alpine 响应式，直接写 DOM + rAF 批处理 ----
    _logBuffer: [],       // 本帧待写入的新行
    _logFlushPending: false,
    LOG_DOM_CAP: 5000,    // DOM 中保留的最大行数
    _enqueueLog(line) {
      this._logBuffer.push(line)
      // 用非响应式计数器代替 currentLogs.push，避免每行触发 Alpine proxy
      this._logCount++
      if (this._logFlushPending) return
      this._logFlushPending = true
      requestAnimationFrame(() => { this._logFlushPending = false; this._flushLogs() })
    },
    _flushLogs() {
      const pre = document.getElementById('log-pre')
      if (!pre) {
        // 组件尚未挂载（x-if 还没渲染），下一帧重试，保留缓冲
        this._logFlushPending = true
        requestAnimationFrame(() => { this._logFlushPending = false; this._flushLogs() })
        return
      }
      if (this._logBuffer.length) {
        // 单次 DOM 写入，合并本帧所有新行
        pre.appendChild(document.createTextNode(this._logBuffer.join('\n') + '\n'))
        this._logBuffer.length = 0
      }
      // 超限时裁剪：展平文本再切片，避免无限增长
      // 粗略估算：按行切，保留末尾 LOG_DOM_CAP 行
      if (pre.childNodes.length > 32) {
        const text = pre.textContent
        const lines = text.split('\n')
        if (lines.length > this.LOG_DOM_CAP) {
          pre.textContent = lines.slice(lines.length - this.LOG_DOM_CAP).join('\n')
        } else if (pre.childNodes.length > 1) {
          // 将多个 textNode 合并为一个，降低后续追加开销
          pre.textContent = text
        }
      }
      // 接近底部才自动跟随
      const area = document.getElementById('log-area')
      if (area) {
        const nearBottom = area.scrollHeight - area.scrollTop - area.clientHeight < 120
        if (nearBottom) area.scrollTop = area.scrollHeight
      }
    },
    _resetLogPre() {
      this._logBuffer.length = 0
      this._logFlushPending = false
      const pre = document.getElementById('log-pre')
      if (pre) pre.textContent = ''
    },
    _setLogPreFromArray(lines) {
      this._logBuffer.length = 0
      const pre = document.getElementById('log-pre')
      if (pre) pre.textContent = (lines || []).join('\n')
    },

    // 兼容保留：手动“滚动到底部”按钮
    scheduleScrollLogsBottom() {
      const area = document.getElementById('log-area')
      if (!area) return
      const nearBottom = area.scrollHeight - area.scrollTop - area.clientHeight < 120
      if (nearBottom) area.scrollTop = area.scrollHeight
    },

    stopSSE() { if (this.sseSource) { this.sseSource.close(); this.sseSource = null } },

    async refreshDetail() {
      if (!this.currentRun) return
      const r = await fetch(`/api/runs/${this.currentRun.run_id}`)
      if (r.ok) {
        const data = await r.json()
        this.currentRun.status = data.status; this.currentRun.ended_at = data.ended_at; this.currentRun.error = data.error; this.currentLogs = data.logs ?? []
        this._setLogPreFromArray(this.currentLogs)
      }
    },

    async loadReport() {
      if (!this.currentRun) return
      const r = await fetch(`/api/runs/${this.currentRun.run_id}/report`)
      if (r.ok) {
        this.currentReport = await r.json()
        this.$nextTick(() => { this.renderChart(); this.renderRadar() })
      }
    },

    openHtmlReport() {
      if (!this.currentRun) return
      window.open(`/api/runs/${this.currentRun.run_id}/report-html`, '_blank')
    },

    renderChart() {
      const canvas = document.getElementById('modelChart')
      if (!canvas || !this.currentReport) return
      Chart.getChart(canvas)?.destroy()
      const dims = this.currentReport.dimensions?.by_model ?? []
      if (!dims.length) return
      const labels = dims.map(d => d.model)
      new Chart(canvas, {
        type: 'bar',
        data: {
          labels,
          datasets: [
            { label:'编译',  data: dims.map(d => +((d.compile_pass_rate||0)*100).toFixed(1)), backgroundColor:'rgba(14,165,233,.7)' },
            { label:'样测',  data: dims.map(d => +((d.avg_test_pass_rate||0)*100).toFixed(1)), backgroundColor:'rgba(16,185,129,.7)' },
            { label:'覆盖率',data: dims.map(d => +((d.avg_line_coverage||0)*100).toFixed(1)), backgroundColor:'rgba(245,158,11,.7)' },
            { label:'变异',  data: dims.map(d => +((d.avg_mutation_score||0)*100).toFixed(1)), backgroundColor:'rgba(239,68,68,.7)' },
          ]
        },
        options: {
          responsive: true, maintainAspectRatio: false,
          scales: {
            y: { min:0, max:100, ticks:{ color:'#334155', font:{size:11} }, grid:{ color:'#0f172a' } },
            x: { ticks:{ color:'#64748b', font:{size:11} }, grid:{ display:false } }
          },
          plugins: { legend: { labels:{ color:'#64748b', font:{size:11}, boxWidth:10, padding:15 } } }
        }
      })
    },

    get reportCards() {
      const s = this.currentReport?.summary ?? {}
      return [
        { label:'样本数',  value: s.total_samples ?? 0, color:'kpi-total' },
        { label:'编译通过率',  value: pct(s.compile_pass_rate), color: pctColor(s.compile_pass_rate) },
        { label:'样本测试通过率',  value: pct(s.sample_test_pass_rate ?? s.test_pass_rate), color: pctColor(s.sample_test_pass_rate ?? s.test_pass_rate) },
        { label:'平均行覆盖率', value: pct(s.avg_line_coverage), color: pctColor(s.avg_line_coverage) },
        { label:'平均变异分', value: pct(s.avg_mutation_score), color: pctColor(s.avg_mutation_score) },
      ]
    },

    scrollLogsBottom() {
      const area = document.getElementById('log-area')
      if (area) area.scrollTop = area.scrollHeight
    },

    fmtTime(t) {
      if (!t) return '—'
      return new Date(t).toLocaleString('zh-CN', { month:'short', day:'numeric', hour:'2-digit', minute:'2-digit' })
    },
    formatDate(t) {
      return this.fmtTime(t)
    },
    duration(start, end) {
      if (!start) return '—'
      const ms = (end ? new Date(end) : new Date()) - new Date(start)
      if (ms < 1000) return ms + 'ms'
      if (ms < 60000) return (ms/1000).toFixed(1) + 's'
      return Math.floor(ms/60000) + 'm ' + Math.floor((ms%60000)/1000) + 's'
    },
    pct(v) { return pct(v) },
    pctColor(v) { return pctColor(v) },
    metricPct(v) {
      if (v == null) return '—'
      return (v > 1 ? v : v * 100).toFixed(1) + '%'
    },
    fmtInt(v) { return (v==null || v===0) ? '—' : Math.round(v).toLocaleString('zh-CN') },
    fmtMs(v) {
      if (v==null || v===0) return '—'
      if (v < 1000) return Math.round(v) + 'ms'
      if (v < 60000) return (v/1000).toFixed(1) + 's'
      return (v/60000).toFixed(1) + 'm'
    },
    fmtBytes(v) {
      if (v == null) return '—'
      if (v < 1024) return v + ' B'
      if (v < 1024 * 1024) return (v / 1024).toFixed(1) + ' KB'
      return (v / 1024 / 1024).toFixed(1) + ' MB'
    },

    renderRadar() {
      const canvas = document.getElementById('radarChart')
      if (!canvas || !this.currentReport) return
      Chart.getChart(canvas)?.destroy()
      const dims = this.currentReport.dimensions?.by_model ?? []
      if (!dims.length) return
      // 效率维度 = 1 - 归一化 tokens_per_pass（无数据为 0）
      const tokArr = dims.map(d => d.tokens_per_pass || 0).filter(v => v > 0)
      const maxTok = tokArr.length ? Math.max(...tokArr) : 1
      const palette = [
        { border:'#0ea5e9', bg:'rgba(14,165,233,.15)' },
        { border:'#10b981', bg:'rgba(16,185,129,.15)' },
        { border:'#f59e0b', bg:'rgba(245,158,11,.15)' },
        { border:'#ef4444', bg:'rgba(239,68,68,.15)' },
        { border:'#a855f7', bg:'rgba(168,85,247,.15)' },
      ]
      new Chart(canvas, {
        type: 'radar',
        data: {
          labels: ['编译','测试','覆盖率','变异','效率'],
          datasets: dims.map((d, i) => ({
            label: d.model,
            data: [
              +((d.compile_pass_rate||0)*100).toFixed(1),
              +((d.avg_test_pass_rate||0)*100).toFixed(1),
              +((d.avg_line_coverage||0)*100).toFixed(1),
              +((d.avg_mutation_score||0)*100).toFixed(1),
              d.tokens_per_pass > 0 ? +((1 - d.tokens_per_pass/maxTok)*100).toFixed(1) : 0,
            ],
            borderColor: palette[i % palette.length].border,
            backgroundColor: palette[i % palette.length].bg,
            borderWidth: 2,
            pointRadius: 3,
            pointBackgroundColor: palette[i % palette.length].border,
          }))
        },
        options: {
          responsive: true, maintainAspectRatio: false,
          scales: {
            r: {
              min: 0, max: 100,
              ticks: { stepSize: 25, color: '#475569', font:{size:10}, backdropColor: 'transparent' },
              grid: { color: '#1e293b' },
              angleLines: { color: '#1e293b' },
              pointLabels: { color: '#cbd5e1', font: { size: 12, weight: '500' } }
            }
          },
          plugins: { legend: { labels: { color:'#94a3b8', font:{size:11}, boxWidth:12, padding:12 } } }
        }
      })
    },

    get heatmap() {
      const rep = this.currentReport
      if (!rep) return { cols: [], rows: [] }
      const byMS = rep.by_model_scenario ?? rep.dimensions?.by_model_scenario ?? []
      if (!byMS.length) return { cols: [], rows: [] }
      const colSet = new Map()
      const grid = new Map()
      for (const r of byMS) {
        const col = `${r.language}·${r.scenario}`
        colSet.set(col, (colSet.get(col) || 0) + 1)
        if (!grid.has(r.model)) grid.set(r.model, {})
        grid.get(r.model)[col] = r.avg_mutation_score
      }
      const cols = [...colSet.keys()].sort()
      const rows = [...grid.entries()].sort((a,b)=>a[0].localeCompare(b[0])).map(([model, m]) => ({
        model,
        values: cols.map(c => (c in m) ? m[c] : null),
      }))
      return { cols, rows }
    },

    heatCell(v) {
      if (v == null) return 'background:var(--heat-null);color:var(--fg-subtle)'
      // 色阶：0 → 暗红，0.5 → 暗黄，1 → 暖橙绿
      const clamped = Math.max(0, Math.min(1, v))
      const hue = 10 + clamped * 130 // 10(红) → 140(绿)
      const alpha = 0.15 + clamped * 0.45
      const fg = clamped > 0.55 ? '#0b1120' : '#e2e8f0'
      return `background:hsla(${hue},70%,50%,${alpha});color:${fg};font-weight:600`
    },
  }
}

function pct(v) { return v==null ? '—' : (v*100).toFixed(1)+'%' }
function pctColor(v) {
  if (v==null) return 'pct-null'
  if (v >= .8) return 'pct-good'
  if (v >= .6) return 'pct-warn'
  return 'pct-bad'
}

