function app() {
  return {
    page: 'dashboard',
    nav: [
      { id:'dashboard', icon:'<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="7" height="7" rx="1"/><rect x="14" y="3" width="7" height="7" rx="1"/><rect x="3" y="14" width="7" height="7" rx="1"/><rect x="14" y="14" width="7" height="7" rx="1"/></svg>', label:'总览' },
      { id:'new-run',   icon:'<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="5 3 19 12 5 21 5 3"/></svg>', label:'新建任务' },
      { id:'runs',      icon:'<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="8" y1="6" x2="21" y2="6"/><line x1="8" y1="12" x2="21" y2="12"/><line x1="8" y1="18" x2="21" y2="18"/><line x1="3" y1="6" x2="3.01" y2="6"/><line x1="3" y1="12" x2="3.01" y2="12"/><line x1="3" y1="18" x2="3.01" y2="18"/></svg>', label:'任务列表' },
      { id:'agents',    icon:'<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 3l7 4v5c0 4.5-3 7.5-7 9-4-1.5-7-4.5-7-9V7l7-4z"/><path d="M9 12h6"/><path d="M12 9v6"/></svg>', label:'Agent 接入' },
      { id:'database',  icon:'<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><ellipse cx="12" cy="5" rx="8" ry="3"/><path d="M4 5v6c0 1.7 3.6 3 8 3s8-1.3 8-3V5"/><path d="M4 11v6c0 1.7 3.6 3 8 3s8-1.3 8-3v-6"/></svg>', label:'资产管理' },
      { id:'environment', icon:'<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/><path d="M9 12l2 2 4-5"/></svg>', label:'环境检查' },
      { id:'models',    icon:'<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M12 1v6m0 10v6m11-11h-6m-10 0H1m15.5-6.5l-4.25 4.25M7.75 16.25L3.5 20.5M20.5 20.5l-4.25-4.25M7.75 7.75L3.5 3.5"/></svg>', label:'模型管理' },
    ],
    config: null,
    runs: [],
    runsLoading: false,
    runFilter: '',
    statusFilter: '',
    agentSubjectQuery: '',
    agentWizardOpen: false,
    agentSaving: false,
    agentFormError: '',
    agentPreset: 'generic',
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
    // 数据库管理
    dbTab: 'overview', // overview | runs | evaluation | assets
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
    // 新增：筛选参数
    dbGenCaseFilter: { run_id:'', model:'', language:'' },
    dbEvalRunFilter: { run_id:'' },
    dbStageFilter: { evaluation_run_id:'' },
    dbSampleFilter: { language:'', class:'' },
    dbAssetFilter: { subject:'', language:'', sample:'' },
    dbReportFilter: { run_id:'' },
    dbRunArtifactFilter: { run_id:'' },
    form: {
      run_id:'', models:[], subjects:[], combinations:[{framework:'model_api',model:'deepseek-v4-flash',skill:'no_skill'}], languages:[], class:'self_contained', scenario:'', level:'',
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
      return this.env.eval_image_name || 'utbench:latest'
    },
    get activeBuildDockerfile() {
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
    _logCount: 0,
    _comboKey: 0,
    partialsLoaded: false,

    async init() {
      this.applyTheme()
      await this.loadPartials()
      await this.loadConfig()
      await this.loadModels()
      await this.loadEnv()
      await this.loadRuns()
      await this.loadDatabase()
      this._startTimers()
      this.$nextTick(() => this._initComboDelegation())
    },
    async loadPartials() {
      const nodes = Array.from(document.querySelectorAll('[data-partial]'))
      await Promise.all(nodes.map(async node => {
        const name = node.getAttribute('data-partial')
        if (!name) return
        const r = await fetch(`/partials/${name}.html`, { cache: 'no-store' })
        if (!r.ok) throw new Error(`load partial ${name}: HTTP ${r.status}`)
        if (window.Alpine?.destroyTree && node._partialInitialized) {
          window.Alpine.destroyTree(node)
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
    },
    _stopTimers() {
      if (this._timerRuns) { clearInterval(this._timerRuns); this._timerRuns = null }
      if (this._timerDb) { clearInterval(this._timerDb); this._timerDb = null }
      if (this._timerEnv) { clearInterval(this._timerEnv); this._timerEnv = null }
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
        if (!prev && this.env.docker_available && this.env.image_present) {
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

    openBuildImage(target = 'eval') { this.buildTarget = target || 'eval'; this.buildModalOpen = true; this.reattachBuild() },
    closeBuildModal() { this.buildModalOpen = false; this.stopBuildSSE() },

    buildStatusColor(s) {
      const map = { running: 'color:var(--accent)', pending: 'color:var(--yellow)', completed: 'color:var(--green)', failed: 'color:var(--red)', canceled: 'color:var(--fg-muted)' }
      return map[s] || 'color:var(--fg-muted)'
    },

    defaultBuildTarget() {
      return 'eval'
    },

    buildTargetText(target) {
      return '评测镜像'
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

    async reattachBuild() {
      try {
        const r = await fetch('/api/env/build-image')
        const data = await r.json()
        if (data && data.build_id) {
          this.buildId = data.build_id
          this.buildTarget = data.target || this.buildTarget || 'eval'
          this.buildStatus = data.status
          this.buildError = data.error || ''
          const detail = await fetch('/api/env/build-image/' + this.buildId)
          if (detail.ok) { const d = await detail.json(); this.buildLogs = d.logs || [] }
          if (this.buildStatus === 'running' || this.buildStatus === 'pending') { this.startBuildSSE(this.buildId) }
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
      this.stopSSE(); this.page = id
      if (id === 'models') this.loadModels()
      if (id === 'agents') { this.loadAPIKeys(); this.loadModels(); this.loadEnv() }
      if (id === 'database') this.loadDBTabData()
      if (id === 'environment') {
        if (!this.environment) this.loadEnvironment()
        this.loadAPIKeys()
      }
    },
    get pageTitle() {
      const map = { dashboard:'总览', 'new-run':'新建任务', runs:'任务列表', agents:'Agent 接入', database:'资产管理', 'run-detail':'任务详情', environment:'环境检查', models:'模型管理' }
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
    openAgentWizard() {
      this.agentFormError = ''
      this.agentWizardOpen = true
      if (!this.agentForm.command) this.applyAgentPreset(this.agentPreset || 'generic')
    },
    closeAgentWizard() { this.agentWizardOpen = false },
    applyAgentPreset(preset) {
      this.agentPreset = preset
      const base = {
        image: 'utbench-agent-base:latest',
        timeout_seconds: 600,
        network_disabled: false,
        compatible_languages: ['python', 'go', 'java', 'cpp'],
        output_globs: 'generated_test.py\ntest_*.py\n*_test.go\n*Test.java\n*test*.cpp',
      }
      const presets = {
        generic: {
          name: this.agentForm.name || 'my-agent',
          env_from_host: this.agentForm.env_from_host || 'MY_AGENT_API_KEY',
          command: 'PROMPT="$(cat {{.ContainerPrompt}})" && my-agent --model "{{.ModelID}}" --prompt "$PROMPT"',
        },
        claudecode: {
          name: this.agentForm.name || 'claudecode-custom',
          env_from_host: 'ANTHROPIC_API_KEY\nANTHROPIC_AUTH_TOKEN\nANTHROPIC_BASE_URL',
          command: 'mkdir -p "{{.ContainerWorkdir}}/.claude" &&\nPROMPT="$(cat {{.ContainerPrompt}})" &&\nclaude -p --output-format stream-json --verbose --permission-mode bypassPermissions "$PROMPT"',
        },
        opencode: {
          name: this.agentForm.name || 'opencode-custom',
          env_from_host: 'DEEPSEEK_API_KEY\nDASHSCOPE_API_KEY\nMINIMAX_API_KEY\nARK_API_KEY\nBIGMODEL_API_KEY',
          command: 'PROMPT="$(cat {{.ContainerPrompt}})" &&\nopencode run --print-logs --dangerously-skip-permissions --model "{{.ModelID}}" "$PROMPT"',
        },
      }
      const next = { ...base, ...(presets[preset] || presets.generic) }
      this.agentForm = { ...this.agentForm, ...next }
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
      const f = this.agentForm
      if (!f.name.trim()) { this.agentFormError = '请填写 Agent 名称'; return }
      if (!f.command.trim()) { this.agentFormError = '请填写启动命令'; return }
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
        this.showToast('Agent 已添加：' + data.name, 'ok')
        this.agentWizardOpen = false
        await this.loadConfig()
      } catch (e) {
        this.agentFormError = e.message || String(e)
        this.showToast('添加 Agent 失败：' + this.agentFormError, 'err', 6000)
      } finally {
        this.agentSaving = false
      }
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
        this.showToast('已创建 Docker 重跑任务 ' + data.run_id, 'ok')
        await this.loadRuns()
        // 切到任务列表方便查看新任务状态
        this.stopSSE(); this.page = 'runs'
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
    renderSelectFw(idx, current) {
      const opts = this._comboFrameworks().map(fw =>
        `<option value="${this._escapeHtml(fw.value)}" ${fw.value===current?'selected':''}>${this._escapeHtml(fw.label)}</option>`
      ).join('')
      return `<select class="input-base text-[12px]" data-combo-idx="${idx}" data-combo-type="fw">${opts}</select>`
    },
    renderSelectModel(idx, framework, current) {
      const models = this._comboModels(framework)
      if (!models.includes(current) && models.length) current = models[0]
      const modelMap = {}
      ;this._modelRows().forEach(m => { modelMap[m.name] = m })
      const opts = models.map(m => {
        const info = modelMap[m]
        const keyOk = info?.api_key_set
        const marker = keyOk ? ' ✓' : ' ✗'
        return `<option value="${this._escapeHtml(m)}" ${m===current?'selected':''}>${this._escapeHtml(m)}${marker}</option>`
      }).join('')
      return `<select class="input-base text-[12px]" data-combo-idx="${idx}" data-combo-type="model">${opts}</select>`
    },
    renderSelectSkill(idx, framework, current) {
      const skills = this._comboSkills(framework)
      if (!skills.includes(current)) current = 'no_skill'
      const opts = skills.map(s =>
        `<option value="${this._escapeHtml(s)}" ${s===current?'selected':''}>${this._escapeHtml(s)}</option>`
      ).join('')
      return `<select class="input-base text-[12px]" data-combo-idx="${idx}" data-combo-type="skill">${opts}</select>`
    },
    _initComboDelegation() {
      const container = document.getElementById('combo-rows')
      if (!container || container._delegated) return
      container._delegated = true
      container.addEventListener('change', (e) => {
        const sel = e.target.closest('select[data-combo-idx]')
        if (!sel) return
        const idx = parseInt(sel.dataset.comboIdx, 10)
        const type = sel.dataset.comboType
        const val = sel.value
        if (type === 'fw') { this.form.combinations[idx].framework = val; this.onCombinationFrameworkChange(idx) }
        else if (type === 'model') { this.form.combinations[idx].model = val }
        else if (type === 'skill') { this.form.combinations[idx].skill = val }
      })
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
      const models = this._comboModels('model_api')
      const defaultModel = models.includes('deepseek-v4-flash') ? 'deepseek-v4-flash' : (models[0] || '')
      this.form.combinations.push({ framework: 'model_api', model: defaultModel, skill: 'no_skill' })
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

    async submitRun() {
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
      try {
        const r = await fetch('/api/runs', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(this.form) })
        const data = await r.json()
        if (!r.ok) { this.formError = data.error || '提交失败'; return }
        await this.loadRuns()
        await this.openRun(data.run_id)
      } catch(e) { this.formError = String(e) }
      finally { this.formSubmitting = false }
    },

    async openRun(runId) {
      this.stopSSE(); this.currentReport = null; this.currentLogs = []; this._logCount = 0; this.detailTab = 'logs'
      this._resetLogPre()
      this.page = 'run-detail'
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
