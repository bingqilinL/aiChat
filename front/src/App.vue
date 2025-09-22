<template>
  <div class="app-container">
    <!-- 左侧历史会话栏 -->
    <div :class="{ 'sidebar-collapsed': !showHistory }" class="sidebar">
      <div class="sidebar-header">
        <div class="header-content">
          <button class="new-chat-btn" @click="startNewChat">
            <svg fill="none" height="16" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" width="16">
              <path d="M12 5v14M5 12h14"/>
            </svg>
            新对话
          </button>
        </div>
        <button class="collapse-btn" @click="toggleHistory">
          <svg fill="none" height="16" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" width="16">
            <path d="M15 18l-6-6 6-6"/>
          </svg>
        </button>
      </div>

      <div class="history-section">
        <div class="section-title">历史对话</div>
        <div class="history-list">
          <div
              v-for="session in historySessions"
              :key="session.id"
              :class="{ 'active': currentSessionId === session.id }"
              class="history-item"
              @click="loadHistory(session.id)"
          >
            <div class="history-content">
              <div class="history-title">{{ session.title }}</div>
              <div class="history-time">{{ formatTime(session.timestamp) }}</div>
            </div>
            <button class="delete-btn" @click.stop="deleteSession(session.id)">
              <svg fill="none" height="14" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" width="14">
                <path d="M3 6h18M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6"/>
              </svg>
            </button>
          </div>
        </div>
        <!-- 清空历史按钮移动到历史栏最下面 -->
        <div class="history-footer">
          <button class="clear-history-btn" @click="clearAllHistory">
            <svg fill="none" height="16" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" width="16">
              <path d="M3 6h18M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6"/>
            </svg>
            清空历史
          </button>
        </div>
      </div>
    </div>

    <!-- 主界面 -->
    <div class="main-content">
      <!-- 顶部工具栏 -->
      <div class="toolbar">
        <button class="menu-btn" @click="toggleHistory">
          <svg fill="none" height="20" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" width="20">
            <path d="M3 12h18M3 6h18M3 18h18"/>
          </svg>
        </button>
        <div class="toolbar-title">AI助手</div>
        <div class="toolbar-actions">
          <input ref="fileInput" class="hidden-file-input" type="file" @change="onFileChange" />
          <button class="action-btn" title="加载知识库" @click="triggerUpload">
            <svg fill="none" height="18" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" width="18">
              <path d="M12 5v14M5 12h14"/>
            </svg>
            加载知识库
          </button>
          <!-- 用户操作按钮 -->
          <div class="user-action-container">
            <button class="action-btn" title="用户操作" @click="toggleUserMenu">
              <svg fill="none" height="18" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" width="18">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                <circle cx="12" cy="7" r="4"/>
              </svg>
              {{ isLoggedIn ? userInfo.username : '登录' }}
            </button>
            
            <!-- 用户菜单下拉框 -->
            <div v-if="showUserMenu" class="user-menu" @click.stop>
              <div v-if="!isLoggedIn" class="user-menu-content">
                <button class="menu-item" @click="showLoginForm = true; showUserMenu = false">
                  <svg fill="none" height="16" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" width="16">
                    <path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/>
                    <polyline points="10,17 15,12 10,7"/>
                    <line x1="15" x2="3" y1="12" y2="12"/>
                  </svg>
                  登录
                </button>
                <button class="menu-item" @click="showRegisterForm = true; showUserMenu = false">
                  <svg fill="none" height="16" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" width="16">
                    <path d="M16 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                    <circle cx="8.5" cy="7" r="4"/>
                    <line x1="20" x2="20" y1="8" y2="14"/>
                    <line x1="23" x2="17" y1="11" y2="11"/>
                  </svg>
                  注册
                </button>
              </div>
              <div v-else class="user-menu-content">
                <div class="user-info">
                  <div class="username">{{ userInfo.username }}</div>
                  <div class="user-email">{{ userInfo.email || '未设置邮箱' }}</div>
                </div>
                <button class="menu-item logout" @click="logout">
                  <svg fill="none" height="16" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" width="16">
                    <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
                    <polyline points="16,17 21,12 16,7"/>
                    <line x1="21" x2="9" y1="12" y2="12"/>
                  </svg>
                  退出登录
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>


      <!-- 登录表单模态框 -->
      <div v-if="showLoginForm" class="modal-overlay" @mousedown="handleModalOverlayClick" @click.prevent>
        <div class="modal-content" @mousedown.stop @click.stop>
          <div class="modal-header">
            <h3>用户登录</h3>
            <button class="close-btn" @click="showLoginForm = false">
              <svg fill="none" height="20" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" width="20">
                <line x1="18" x2="6" y1="6" y2="18"/>
                <line x1="6" x2="18" y1="6" y2="18"/>
              </svg>
            </button>
          </div>
          <form class="auth-form" @submit.prevent="login">
            <div class="form-group">
              <label>用户名</label>
              <input v-model="loginForm.username" placeholder="请输入用户名" required type="text">
            </div>
            <div class="form-group">
              <label>密码</label>
              <input v-model="loginForm.password" placeholder="请输入密码" required type="password">
            </div>
            <button :disabled="loginLoading" class="submit-btn" type="submit">
              {{ loginLoading ? '登录中...' : '登录' }}
            </button>
          </form>
          <!-- 登录消息显示 -->
          <div v-if="loginMessage" class="auth-message error">
            {{ loginMessage }}
          </div>
        </div>
      </div>

      <!-- 注册表单模态框 -->
      <div v-if="showRegisterForm" class="modal-overlay" @mousedown="handleModalOverlayClick" @click.prevent>
        <div class="modal-content" @mousedown.stop @click.stop>
          <div class="modal-header">
            <h3>用户注册</h3>
            <button class="close-btn" @click="showRegisterForm = false">
              <svg fill="none" height="20" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" width="20">
                <line x1="18" x2="6" y1="6" y2="18"/>
                <line x1="6" x2="18" y1="6" y2="18"/>
              </svg>
            </button>
          </div>
          <form class="auth-form" @submit.prevent="register">
            <div class="form-group">
              <label>用户名</label>
              <input v-model="registerForm.username" placeholder="请输入用户名" required type="text">
            </div>
            <div class="form-group">
              <label>密码</label>
              <input v-model="registerForm.password" placeholder="请输入密码" required type="password">
            </div>
            <div class="form-group">
              <label>确认密码</label>
              <input v-model="registerForm.confirmPassword" placeholder="请再次输入密码" required type="password">
            </div>
            <div class="form-group">
              <label>邮箱（可选）</label>
              <input v-model="registerForm.email" placeholder="请输入邮箱" type="email">
            </div>
            <button :disabled="registerLoading" class="submit-btn" type="submit">
              {{ registerLoading ? '注册中...' : '注册' }}
            </button>
          </form>
          <!-- 注册消息显示 -->
          <div v-if="registerMessage" class="auth-message error">
            {{ registerMessage }}
          </div>
        </div>
      </div>

      <!-- 登出消息显示 -->
      <div v-if="showLogoutMessage" class="logout-message-overlay">
        <div class="logout-message">
          {{ logoutMessage }}
        </div>
      </div>

      <!-- 聊天区域 -->
      <div class="chat-container">
        <div ref="messagesContainer" class="chat-messages">
          <!-- 欢迎消息 -->
          <div v-if="messages.length === 0" class="welcome-section">
            <div class="welcome-content">
              <h2>欢迎使用AI助手</h2>
              <p>我可以帮助您解决学习、工作和生活中的各种问题</p>
              <div class="suggestions">
                <div class="suggestion-item" @click="sendSuggestion('请帮我写一个Python函数')">
                  写代码
                </div>
                <div class="suggestion-item" @click="sendSuggestion('解释一下什么是机器学习')">
                  学习知识
                </div>
                <div class="suggestion-item" @click="sendSuggestion('帮我制定一个学习计划')">
                  制定计划
                </div>
                <div class="suggestion-item" @click="sendSuggestion('翻译这段文字')">
                  翻译助手
                </div>
              </div>
            </div>
          </div>

          <!-- 消息列表 -->
          <div
              v-for="(message, index) in messages"
              :key="message.id"
              :class="message.sender"
              class="message-wrapper"
          >
            <div class="message-container">
              <div class="message-avatar">
                <div v-if="message.sender === 'user'" class="user-avatar">
                  <svg fill="none" height="20" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" width="20">
                    <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                    <circle cx="12" cy="7" r="4"/>
                  </svg>
                </div>
                <div v-else class="ai-avatar">
                  <svg fill="none" height="20" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" width="20">
                    <path d="M12 2L2 7l10 5 10-5-10-5z"/>
                    <path d="M2 17l10 5 10-5"/>
                    <path d="M2 12l10 5 10-5"/>
                  </svg>
                </div>
              </div>
              <div class="message-content">
                <div class="message-text" v-html="renderMarkdown(message.content)"></div>
                <div class="message-time">{{ formatMessageTime(message.timestamp) }}</div>
              </div>
            </div>
          </div>

          <!-- 流式输出显示区域 -->
          <div v-if="streamingMessage" class="message-wrapper ai">
            <div class="message-container">
              <div class="message-avatar">
                <div class="ai-avatar">
                  <svg fill="none" height="20" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" width="20">
                    <path d="M12 2L2 7l10 5 10-5-10-5z"/>
                    <path d="M2 17l10 5 10-5"/>
                    <path d="M2 12l10 5 10-5"/>
                  </svg>
                </div>
              </div>
              <div class="message-content">
                <div class="message-text" v-html="renderMarkdown(streamingContent)"></div>
                <div class="typing-indicator">
                  <span></span>
                  <span></span>
                  <span></span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 输入区域 -->
        <div class="input-container">
          <div class="input-wrapper">
            <textarea
                ref="inputRef"
                v-model="inputMessage"
                :disabled="loading"
                placeholder="输入您的问题..."
                rows="1"
                @keydown="handleKeyDown"
            ></textarea>
            <button
                :disabled="loading || !inputMessage.trim()"
                class="send-btn"
                @click="sendMessage"
            >
              <svg v-if="!loading" fill="none" height="20" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" width="20">
                <path d="M22 2L11 13M22 2l-7 20-4-9-9-4 20-7z"/>
              </svg>
              <div v-else class="loading-spinner"></div>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted, nextTick, watch } from 'vue'

export default {
  name: 'ChatApp',
  setup() {
    const messages = ref([])
    const inputMessage = ref('')
    const loading = ref(false)
    const messagesContainer = ref(null)
    const inputRef = ref(null)
    const showHistory = ref(true)
    const historySessions = ref([])
    const currentSessionId = ref(null)
    const fileInput = ref(null)

    // 流式输出相关
    const streamingMessage = ref(false)
    const streamingContent = ref('')
    const streamingMessageId = ref('')

    // 用户相关状态
    const isLoggedIn = ref(false)
    const userInfo = ref({})
    const showUserMenu = ref(false)
    const showLoginForm = ref(false)
    const showRegisterForm = ref(false)
    const loginLoading = ref(false)
    const registerLoading = ref(false)
    
    // 消息状态管理
    const loginMessage = ref('')
    const registerMessage = ref('')
    const logoutMessage = ref('')
    const showLogoutMessage = ref(false)
    
    // 登录表单
    const loginForm = ref({
      username: '',
      password: ''
    })
    
    // 注册表单
    const registerForm = ref({
      username: '',
      password: '',
      confirmPassword: '',
      email: ''
    })

    // 增强的Markdown渲染器
    const renderMarkdown = (text) => {
      if (!text) return ''

      // 转义HTML
      let html = text
          .replace(/&/g, '&amp;')
          .replace(/</g, '&lt;')
          .replace(/>/g, '&gt;')

      // 代码块 - 支持语言标识
      html = html.replace(/```(\w+)?\n?([\s\S]*?)```/g, (match, lang, code) => {
        const language = lang || 'text'
        const cleanCode = code.trim()
        return `<pre class="code-block"><code class="language-${language}">${cleanCode}</code></pre>`
      })

      // 行内代码
      html = html.replace(/`([^`]+)`/g, '<code class="inline-code">$1</code>')

      // 粗体
      html = html.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')

      // 斜体
      html = html.replace(/\*(.*?)\*/g, '<em>$1</em>')

      // 链接
      html = html.replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" target="_blank" rel="noopener noreferrer">$1</a>')

      // 标题
      html = html.replace(/^### (.*$)/gm, '<h3>$1</h3>')
      html = html.replace(/^## (.*$)/gm, '<h2>$1</h2>')
      html = html.replace(/^# (.*$)/gm, '<h1>$1</h1>')

      // 有序列表
      html = html.replace(/^\d+\. (.*$)/gm, '<li>$1</li>')
      html = html.replace(/(<li>.*<\/li>)/s, '<ol>$1</ol>')

      // 无序列表
      html = html.replace(/^[-*] (.*$)/gm, '<li>$1</li>')
      html = html.replace(/(<li>.*<\/li>)/s, '<ul>$1</ul>')

      // 引用
      html = html.replace(/^> (.*$)/gm, '<blockquote>$1</blockquote>')

      // 水平线
      html = html.replace(/^---$/gm, '<hr>')

      // 表格支持（简单）
      html = html.replace(/\|(.+)\|/g, (match, content) => {
        const cells = content.split('|').map(cell => cell.trim())
        return `<tr>${cells.map(cell => `<td>${cell}</td>`).join('')}</tr>`
      })

      // 换行处理
      html = html.replace(/\n\n/g, '</p><p>')
      html = html.replace(/\n/g, '<br>')

      // 包装段落
      if (!html.startsWith('<')) {
        html = `<p>${html}</p>`
      }

      return html
    }

    // 切换历史记录侧边栏
    const toggleHistory = () => {
      showHistory.value = !showHistory.value
    }

    // 格式化时间
    const formatTime = (timestamp) => {
      const date = new Date(timestamp)
      const now = new Date()
      const diff = now - date

      if (diff < 60000) return '刚刚'
      if (diff < 3600000) return `${Math.floor(diff / 60000)}分钟前`
      if (diff < 86400000) return `${Math.floor(diff / 3600000)}小时前`
      if (diff < 604800000) return `${Math.floor(diff / 86400000)}天前`

      return date.toLocaleDateString('zh-CN')
    }

    // 格式化消息时间
    const formatMessageTime = (timestamp) => {
      const date = new Date(timestamp)
      return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
    }

    // 开始新对话
    const startNewChat = () => {
      currentSessionId.value = null
      messages.value = []
      inputRef.value?.focus()
    }

    // 删除会话
    const deleteSession = (sessionId) => {
      if (confirm('确定要删除这个对话吗？')) {
        historySessions.value = historySessions.value.filter(s => s.id !== sessionId)
        if (currentSessionId.value === sessionId) {
          startNewChat()
        }
        localStorage.setItem('chatHistory', JSON.stringify(historySessions.value))
      }
    }

    // 清空所有历史记录
    const clearAllHistory = () => {
      if (confirm('确定要清空所有历史记录吗？')) {
        historySessions.value = []
        localStorage.removeItem('chatHistory')
        startNewChat()
      }
    }

    // 发送建议
    const sendSuggestion = (suggestion) => {
      inputMessage.value = suggestion
      sendMessage()
    }

    // 保存当前会话到历史记录
    const saveCurrentSession = () => {
      if (messages.value.length <= 0) return

      const sessionTitle = messages.value.find(m => m.sender === 'user')?.content.substring(0, 30) || '新会话'

      const session = {
        id: currentSessionId.value || Date.now().toString(),
        title: sessionTitle,
        timestamp: currentSessionId.value ? historySessions.value.find(s => s.id === currentSessionId.value)?.timestamp : Date.now(),
        messages: messages.value
      }

      // 更新或添加会话
      const existingIndex = historySessions.value.findIndex(s => s.id === session.id)
      if (existingIndex >= 0) {
        historySessions.value[existingIndex] = session
      } else {
        historySessions.value.unshift(session)
        currentSessionId.value = session.id
      }

      // 限制历史记录数量
      if (historySessions.value.length > 50) {
        historySessions.value = historySessions.value.slice(0, 50)
      }

      // 保存到本地存储
      localStorage.setItem('chatHistory', JSON.stringify(historySessions.value))
    }

    // 加载历史记录
    const loadHistory = (sessionId) => {
      const session = historySessions.value.find(s => s.id === sessionId)
      if (session) {
        messages.value = [...session.messages]
        currentSessionId.value = session.id
        scrollToBottom()
      }
    }

    // 滚动到底部
    const scrollToBottom = async () => {
      await nextTick()
      if (messagesContainer.value) {
        messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
      }
    }

    // 处理键盘事件
    const handleKeyDown = (event) => {
      if (event.key === 'Enter' && !event.shiftKey) {
        event.preventDefault()
        sendMessage()
      }
    }

    // 自动调整输入框高度
    const adjustTextareaHeight = () => {
      if (inputRef.value) {
        inputRef.value.style.height = 'auto'
        inputRef.value.style.height = Math.min(inputRef.value.scrollHeight, 120) + 'px'
      }
    }

    // 监听输入内容变化
    watch(inputMessage, adjustTextareaHeight)

    // 处理流式响应
    const processStreamResponse = (response) => {
      return new Promise((resolve, reject) => {
        const reader = response.body.getReader()
        const decoder = new TextDecoder('utf-8')

        // 创建流式消息
        streamingMessage.value = true
        streamingContent.value = ''
        streamingMessageId.value = 'stream_' + Date.now()

        const processStream = async () => {
          try {
            while (true) {
              const { done, value } = await reader.read()
              if (done) break

              const chunk = decoder.decode(value, { stream: true })

              // 处理 SSE 格式的数据
              const lines = chunk.split('\n')
              for (const line of lines) {
                if (line.startsWith('data: ')) {
                  const data = line.substring(6)
                  streamingContent.value += data
                  await scrollToBottom()
                }
              }
            }

            // 流结束，将内容添加到消息列表
            if (streamingContent.value) {
              messages.value.push({
                id: streamingMessageId.value,
                sender: 'ai',
                content: streamingContent.value,
                timestamp: Date.now()
              })
            }

            // 重置流状态
            streamingMessage.value = false
            resolve()
          } catch (error) {
            console.error('Stream reading error:', error)
            streamingMessage.value = false
            reject(error)
          } finally {
            reader.releaseLock()
          }
        }

        processStream()
      })
    }

    // 发送消息
    const sendMessage = async () => {
      const content = inputMessage.value.trim()
      if (!content || loading.value) return

      // 添加用户消息
      const userMessage = {
        id: Date.now().toString(),
        sender: 'user',
        content: content,
        timestamp: Date.now()
      }
      messages.value.push(userMessage)

      // 清空输入框
      inputMessage.value = ''
      loading.value = true
      streamingMessage.value = false

      // 滚动到底部
      await scrollToBottom()

      try {
        // 调用后端API（流式）
        const response = await fetch('http://localhost:9343/AIchat/li', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          credentials: 'include',
          body: JSON.stringify({
            id: Date.now().toString(),
            content: content
          })
        })

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`)
        }

        // 处理流式响应
        await processStreamResponse(response)
      } catch (error) {
        // 错误处理
        const errorMessage = {
          id: Date.now().toString(),
          sender: 'ai',
          content: '抱歉，我遇到了一些问题，请稍后再试。',
          timestamp: Date.now()
        }
        messages.value.push(errorMessage)
        console.error('Error:', error)
      } finally {
        loading.value = false
        await scrollToBottom()
        // 保存会话
        saveCurrentSession()
        // 聚焦输入框
        inputRef.value?.focus()
      }
    }

    // 触发选择文件
    const triggerUpload = () => {
      fileInput.value && fileInput.value.click()
    }

    // 选择文件变更时上传
    const onFileChange = async (e) => {
      const file = e.target.files && e.target.files[0]
      if (!file) return
      try {
        const form = new FormData()
        form.append('file', file)
        const resp = await fetch('http://localhost:9343/AIchat/upload', {
          method: 'POST',
          body: form
        })
        if (!resp.ok) throw new Error('上传失败')
        alert('上传成功并已入库: ' + file.name)
      } catch (err) {
        console.error(err)
        alert('上传失败，请重试')
      } finally {
        e.target.value = ''
      }
    }

    // 用户相关方法
    const toggleUserMenu = () => {
      showUserMenu.value = !showUserMenu.value
    }

    const login = async () => {
      if (!loginForm.value.username || !loginForm.value.password) {
        loginMessage.value = '请填写用户名和密码'
        return
      }

      loginLoading.value = true
      loginMessage.value = ''
      try {
        const response = await fetch('http://localhost:9343/user/login', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            username: loginForm.value.username,
            password: loginForm.value.password
          })
        })

        const result = await response.json()
        if (response.ok) {
          isLoggedIn.value = true
          userInfo.value = result.user
          showLoginForm.value = false
          loginForm.value = { username: '', password: '' }
          loginMessage.value = ''
        } else {
          loginMessage.value = result.error || '登录失败'
        }
      } catch (error) {
        console.error('Login error:', error)
        loginMessage.value = '登录失败，请重试'
      } finally {
        loginLoading.value = false
      }
    }

    const register = async () => {
      if (!registerForm.value.username || !registerForm.value.password) {
        registerMessage.value = '请填写用户名和密码'
        return
      }

      if (registerForm.value.password !== registerForm.value.confirmPassword) {
        registerMessage.value = '两次输入的密码不一致'
        return
      }

      registerLoading.value = true
      registerMessage.value = ''
      try {
        const response = await fetch('http://localhost:9343/user/register', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            username: registerForm.value.username,
            password: registerForm.value.password,
            email: registerForm.value.email
          })
        })

        const result = await response.json()
        if (response.ok) {
          showRegisterForm.value = false
          registerForm.value = { username: '', password: '', confirmPassword: '', email: '' }
          registerMessage.value = ''
          // 注册成功后自动打开登录表单
          showLoginForm.value = true
        } else {
          registerMessage.value = result.error || '注册失败'
        }
      } catch (error) {
        console.error('Register error:', error)
        registerMessage.value = '注册失败，请重试'
      } finally {
        registerLoading.value = false
      }
    }

    const logout = async () => {
      try {
        const response = await fetch('http://localhost:9343/user/logout', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          credentials: 'include'
        })

        if (response.ok) {
          isLoggedIn.value = false
          userInfo.value = {}
          showUserMenu.value = false
          logoutMessage.value = '已退出登录'
          showLogoutMessage.value = true
          // 3秒后自动隐藏消息
          setTimeout(() => {
            showLogoutMessage.value = false
            logoutMessage.value = ''
          }, 3000)
        } else {
          const result = await response.json()
          logoutMessage.value = result.error || '退出登录失败'
          showLogoutMessage.value = true
          setTimeout(() => {
            showLogoutMessage.value = false
            logoutMessage.value = ''
          }, 3000)
        }
      } catch (error) {
        console.error('Logout error:', error)
        logoutMessage.value = '退出登录失败，请重试'
        showLogoutMessage.value = true
        setTimeout(() => {
          showLogoutMessage.value = false
          logoutMessage.value = ''
        }, 3000)
      }
    }

    // 点击外部关闭用户菜单
    const handleClickOutside = (event) => {
      if (showUserMenu.value && !event.target.closest('.user-action-container')) {
        showUserMenu.value = false
      }
    }

    // 处理模态框外部点击
    const handleModalOverlayClick = (event) => {
      // 只有在点击的是模态框背景时才关闭
      if (event.target === event.currentTarget) {
        showLoginForm.value = false
        showRegisterForm.value = false
      }
    }

    // 处理键盘事件
    const handleKeydown = (event) => {
      // ESC键关闭模态框
      if (event.key === 'Escape') {
        showLoginForm.value = false
        showRegisterForm.value = false
        showUserMenu.value = false
      }
    }

    // 初始化
    onMounted(() => {
      // 从本地存储加载历史记录
      const savedHistory = localStorage.getItem('chatHistory')
      if (savedHistory) {
        try {
          historySessions.value = JSON.parse(savedHistory)
        } catch (e) {
          console.error('Failed to parse history', e)
        }
      }

      // 添加点击外部关闭菜单的事件监听
      document.addEventListener('click', handleClickOutside)
      // 添加键盘事件监听
      document.addEventListener('keydown', handleKeydown)

      // 聚焦输入框
      inputRef.value?.focus()
    })

    // 组件卸载时移除事件监听
    onUnmounted(() => {
      document.removeEventListener('click', handleClickOutside)
      document.removeEventListener('keydown', handleKeydown)
    })

    return {
      messages,
      inputMessage,
      loading,
      messagesContainer,
      inputRef,
      showHistory,
      historySessions,
      currentSessionId,
      streamingMessage,
      streamingContent,
      sendMessage,
      toggleHistory,
      startNewChat,
      deleteSession,
      clearAllHistory,
      loadHistory,
      formatTime,
      formatMessageTime,
      renderMarkdown,
      sendSuggestion,
      handleKeyDown,
      fileInput,
      triggerUpload,
      onFileChange,
      // 用户相关
      isLoggedIn,
      userInfo,
      showUserMenu,
      showLoginForm,
      showRegisterForm,
      loginLoading,
      registerLoading,
      loginForm,
      registerForm,
      loginMessage,
      registerMessage,
      logoutMessage,
      showLogoutMessage,
      toggleUserMenu,
      login,
      register,
      logout,
      handleModalOverlayClick,
      handleKeydown
    }
  }
}
</script>

<style scoped>
* {
  box-sizing: border-box;
}

.app-container {
  display: flex;
  height: 100vh;
  background-color: #f7f7f8;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
}

/* 左侧边栏 */
.sidebar {
  width: 10vw;
  background-color: #ffffff;
  border-right: 1px solid #e5e7eb;
  display: flex;
  flex-direction: column;
  transition: all 0.3s ease;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.1);
}

.sidebar-collapsed {
  width: 0;
  overflow: hidden;
}

.sidebar-header {
  padding: 16px;
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-content {
  flex: 1;
}

.new-chat-btn {
  width: 100%;
  padding: 12px 16px;
  background-color: #10a37f;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: background-color 0.2s;
}

.new-chat-btn:hover {
  background-color: #0d8f6f;
}

.collapse-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 8px;
  border-radius: 6px;
  color: #6b7280;
  transition: all 0.2s;
}

.collapse-btn:hover {
  background-color: #f3f4f6;
  color: #374151;
}

.history-section {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.section-title {
  padding: 16px 16px 8px;
  font-size: 12px;
  font-weight: 600;
  color: #6b7280;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.history-list {
  flex: 1;
  overflow-y: auto;
  padding: 0 8px;
}

.history-item {
  display: flex;
  align-items: center;
  padding: 12px;
  margin: 4px 0;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
}

.history-item:hover {
  background-color: #f3f4f6;
}

.history-item.active {
  background-color: #e0f2fe;
  border: 1px solid #0ea5e9;
}

.history-content {
  flex: 1;
  min-width: 0;
}

.history-title {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.history-time {
  font-size: 12px;
  color: #9ca3af;
}

.delete-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  color: #9ca3af;
  opacity: 0;
  transition: all 0.2s;
}

.history-item:hover .delete-btn {
  opacity: 1;
}

.delete-btn:hover {
  background-color: #fee2e2;
  color: #dc2626;
}

/* 历史栏底部清空按钮 */
.history-footer {
  padding: 16px;
  border-top: 1px solid #e5e7eb;
}

.clear-history-btn {
  width: 100%;
  padding: 12px 16px;
  background-color: #fee2e2;
  color: #dc2626;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  transition: background-color 0.2s;
}

.clear-history-btn:hover {
  background-color: #fecaca;
}

/* 主界面 */
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.toolbar {
  height: 60px;
  background-color: #ffffff;
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  align-items: center;
  padding: 0 20px;
  gap: 16px;
  position: relative;
}

.menu-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 8px;
  border-radius: 6px;
  color: #6b7280;
  transition: all 0.2s;
}

.menu-btn:hover {
  background-color: #f3f4f6;
  color: #374151;
}

.toolbar-title {
  flex: 1;
  font-size: 18px;
  font-weight: 600;
  color: #111827;
}

.toolbar-actions {
  display: flex;
  gap: 8px;
}

.hidden-file-input {
  display: none;
}

.action-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 8px;
  border-radius: 6px;
  color: #6b7280;
  transition: all 0.2s;
}

.action-btn:hover {
  background-color: #f3f4f6;
  color: #374151;
}

/* 用户操作容器 */
.user-action-container {
  position: relative;
}

/* 用户菜单 */
.user-menu {
  position: absolute;
  top: 100%;
  right: 0;
  background-color: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 1000;
  min-width: 200px;
}

.user-menu-content {
  padding: 8px 0;
}

.user-info {
  padding: 12px 16px;
  border-bottom: 1px solid #e5e7eb;
  background-color: #f9fafb;
}

.username {
  font-weight: 600;
  color: #111827;
  margin-bottom: 4px;
}

.user-email {
  font-size: 12px;
  color: #6b7280;
}

.menu-item {
  width: 100%;
  padding: 12px 16px;
  background: none;
  border: none;
  text-align: left;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 14px;
  color: #374151;
  transition: background-color 0.2s;
}

.menu-item:hover {
  background-color: #f3f4f6;
}

.menu-item.logout {
  color: #dc2626;
}

.menu-item.logout:hover {
  background-color: #fee2e2;
}

/* 模态框 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}

.modal-content {
  background-color: #ffffff;
  border-radius: 12px;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
  width: 90%;
  max-width: 400px;
  max-height: 90vh;
  overflow-y: auto;
  position: relative;
  z-index: 1;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid #e5e7eb;
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #111827;
}

.close-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  color: #6b7280;
  transition: all 0.2s;
}

.close-btn:hover {
  background-color: #f3f4f6;
  color: #374151;
}

.auth-form {
  padding: 24px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.form-group input {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  font-size: 14px;
  transition: border-color 0.2s, box-shadow 0.2s;
  box-sizing: border-box;
}

.form-group input:focus {
  outline: none;
  border-color: #10a37f;
  box-shadow: 0 0 0 3px rgba(16, 163, 127, 0.1);
}

.submit-btn {
  width: 100%;
  padding: 12px 16px;
  background-color: #10a37f;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;
}

.submit-btn:hover:not(:disabled) {
  background-color: #0d8f6f;
}

.submit-btn:disabled {
  background-color: #d1d5db;
  cursor: not-allowed;
}

/* 认证消息样式 */
.auth-message {
  margin-top: 16px;
  padding: 12px 16px;
  border-radius: 8px;
  font-size: 14px;
  text-align: center;
  animation: slideIn 0.3s ease-out;
}

.auth-message.error {
  background-color: #fee2e2;
  color: #dc2626;
  border: 1px solid #fecaca;
}

.auth-message.success {
  background-color: #d1fae5;
  color: #059669;
  border: 1px solid #a7f3d0;
}

/* 登出消息样式 */
.logout-message-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 1000;
  animation: fadeIn 0.3s ease-out;
}

.logout-message {
  background-color: white;
  padding: 24px 32px;
  border-radius: 12px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
  font-size: 16px;
  font-weight: 500;
  color: #374151;
  text-align: center;
  animation: slideUp 0.3s ease-out;
}

/* 动画效果 */
@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 聊天容器 */
.chat-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 20px 25vw; /* 中心对话框左右留白各25% */
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* 欢迎区域 */
.welcome-section {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.welcome-content {
  text-align: center;
  max-width: 600px;
}

.welcome-content h2 {
  font-size: 32px;
  font-weight: 700;
  color: #111827;
  margin-bottom: 16px;
}

.welcome-content p {
  font-size: 18px;
  color: #6b7280;
  margin-bottom: 32px;
}

.suggestions {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 12px;
}

.suggestion-item {
  padding: 16px;
  background-color: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 14px;
  color: #374151;
}

.suggestion-item:hover {
  background-color: #f9fafb;
  border-color: #10a37f;
  transform: translateY(-2px);
}

/* 消息样式 */
.message-wrapper {
  display: flex;
  width: 100%;
}

.message-wrapper.user {
  justify-content: flex-end;
}

.message-wrapper.ai {
  justify-content: flex-start;
}

.message-container {
  display: flex;
  gap: 12px;
  max-width: 50%;
  align-items: flex-start;
}

.message-wrapper.user .message-container {
  flex-direction: row-reverse;
}

.message-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.user-avatar {
  background-color: #10a37f;
  color: white;
}

.ai-avatar {
  background-color: #f3f4f6;
  color: #6b7280;
}

.message-content {
  flex: 1;
  min-width: 0;
}

.message-text {
  background-color: #ffffff;
  padding: 16px 20px;
  border-radius: 18px;
  font-size: 15px;
  line-height: 1.6;
  color: #374151;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  word-wrap: break-word;
}

.message-wrapper.user .message-text {
  background-color: #10a37f;
  color: white;
  border-bottom-right-radius: 4px;
}

.message-wrapper.ai .message-text {
  border-bottom-left-radius: 4px;
}

.message-time {
  font-size: 12px;
  color: #9ca3af;
  margin-top: 8px;
  text-align: right;
}

.message-wrapper.ai .message-time {
  text-align: left;
}

/* 打字指示器 */
.typing-indicator {
  display: flex;
  gap: 4px;
  margin-top: 8px;
}

.typing-indicator span {
  width: 6px;
  height: 6px;
  background-color: #9ca3af;
  border-radius: 50%;
  animation: typing 1.4s infinite ease-in-out;
}

.typing-indicator span:nth-child(2) {
  animation-delay: 0.2s;
}

.typing-indicator span:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes typing {
  0%, 60%, 100% {
    transform: translateY(0);
    opacity: 0.5;
  }
  30% {
    transform: translateY(-10px);
    opacity: 1;
  }
}

/* 输入区域 */
.input-container {
  padding: 20px;
  background-color: #ffffff;
  border-top: 1px solid #e5e7eb;
}

.input-wrapper {
  max-width: 800px;
  margin: 0 auto;
  display: flex;
  gap: 12px;
  align-items: flex-end;
  background-color: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 24px;
  padding: 12px 16px;
  transition: all 0.2s;
}

.input-wrapper:focus-within {
  border-color: #10a37f;
  box-shadow: 0 0 0 3px rgba(16, 163, 127, 0.1);
}

.input-wrapper textarea {
  flex: 1;
  border: none;
  background: none;
  outline: none;
  resize: none;
  font-size: 15px;
  line-height: 1.5;
  color: #374151;
  min-height: 24px;
  max-height: 120px;
  font-family: inherit;
}

.input-wrapper textarea::placeholder {
  color: #9ca3af;
}

.send-btn {
  width: 40px;
  height: 40px;
  border: none;
  border-radius: 50%;
  background-color: #10a37f;
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  flex-shrink: 0;
}

.send-btn:hover:not(:disabled) {
  background-color: #0d8f6f;
  transform: scale(1.05);
}

.send-btn:disabled {
  background-color: #d1d5db;
  cursor: not-allowed;
  transform: none;
}

.loading-spinner {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top: 2px solid white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* Markdown样式 */
.message-text :deep(p) {
  margin: 8px 0;
  line-height: 1.6;
}

.message-text :deep(.code-block) {
  background-color: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 16px;
  overflow-x: auto;
  margin: 12px 0;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', monospace;
  font-size: 14px;
  line-height: 1.5;
  position: relative;
}

.message-text :deep(.code-block::before) {
  content: attr(data-language);
  position: absolute;
  top: 8px;
  right: 12px;
  font-size: 12px;
  color: #6b7280;
  background-color: #e5e7eb;
  padding: 2px 8px;
  border-radius: 4px;
  text-transform: uppercase;
}

.message-text :deep(.code-block code) {
  background: none;
  padding: 0;
  border-radius: 0;
  font-family: inherit;
  font-size: inherit;
  color: inherit;
}

.message-text :deep(.inline-code) {
  background-color: #f1f3f4;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', monospace;
  font-size: 13px;
  color: #d63384;
}

.message-text :deep(strong) {
  font-weight: 600;
  color: #111827;
}

.message-text :deep(em) {
  font-style: italic;
  color: #374151;
}

.message-text :deep(a) {
  color: #10a37f;
  text-decoration: none;
  border-bottom: 1px solid transparent;
  transition: all 0.2s;
}

.message-text :deep(a:hover) {
  text-decoration: none;
  border-bottom-color: #10a37f;
}

.message-text :deep(h1) {
  font-size: 24px;
  font-weight: 700;
  margin: 20px 0 16px 0;
  color: #111827;
  border-bottom: 2px solid #e5e7eb;
  padding-bottom: 8px;
}

.message-text :deep(h2) {
  font-size: 20px;
  font-weight: 600;
  margin: 18px 0 14px 0;
  color: #111827;
  border-bottom: 1px solid #e5e7eb;
  padding-bottom: 6px;
}

.message-text :deep(h3) {
  font-size: 18px;
  font-weight: 600;
  margin: 16px 0 12px 0;
  color: #111827;
}

.message-text :deep(h4) {
  font-size: 16px;
  font-weight: 600;
  margin: 14px 0 10px 0;
  color: #374151;
}

.message-text :deep(h5) {
  font-size: 14px;
  font-weight: 600;
  margin: 12px 0 8px 0;
  color: #374151;
}

.message-text :deep(h6) {
  font-size: 13px;
  font-weight: 600;
  margin: 10px 0 6px 0;
  color: #6b7280;
}

.message-text :deep(ul) {
  margin: 12px 0;
  padding-left: 20px;
}

.message-text :deep(ol) {
  margin: 12px 0;
  padding-left: 20px;
}

.message-text :deep(li) {
  margin: 4px 0;
  line-height: 1.5;
}

.message-text :deep(blockquote) {
  border-left: 4px solid #10a37f;
  background-color: #f0fdf4;
  padding: 12px 16px;
  margin: 16px 0;
  border-radius: 0 8px 8px 0;
  font-style: italic;
  color: #374151;
}

.message-text :deep(blockquote p) {
  margin: 0;
}

.message-text :deep(hr) {
  border: none;
  height: 1px;
  background-color: #e5e7eb;
  margin: 20px 0;
}

.message-text :deep(table) {
  border-collapse: collapse;
  width: 100%;
  margin: 16px 0;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  overflow: hidden;
}

.message-text :deep(th) {
  background-color: #f9fafb;
  padding: 12px;
  text-align: left;
  font-weight: 600;
  border-bottom: 1px solid #e5e7eb;
}

.message-text :deep(td) {
  padding: 12px;
  border-bottom: 1px solid #f3f4f6;
}

.message-text :deep(tr:last-child td) {
  border-bottom: none;
}

.message-text :deep(tr:nth-child(even)) {
  background-color: #f9fafb;
}

/* 用户消息中的Markdown样式调整 */
.message-wrapper.user .message-text :deep(.code-block) {
  background-color: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.2);
  color: #e2e8f0;
}

.message-wrapper.user .message-text :deep(.inline-code) {
  background-color: rgba(255, 255, 255, 0.2);
  color: #fbb6ce;
}

.message-wrapper.user .message-text :deep(a) {
  color: #a7f3d0;
}

.message-wrapper.user .message-text :deep(a:hover) {
  border-bottom-color: #a7f3d0;
}

.message-wrapper.user .message-text :deep(blockquote) {
  border-left-color: #a7f3d0;
  background-color: rgba(255, 255, 255, 0.1);
  color: #e2e8f0;
}

.message-wrapper.user .message-text :deep(hr) {
  background-color: rgba(255, 255, 255, 0.3);
}

.message-wrapper.user .message-text :deep(table) {
  border-color: rgba(255, 255, 255, 0.2);
}

.message-wrapper.user .message-text :deep(th) {
  background-color: rgba(255, 255, 255, 0.1);
  border-bottom-color: rgba(255, 255, 255, 0.2);
}

.message-wrapper.user .message-text :deep(td) {
  border-bottom-color: rgba(255, 255, 255, 0.1);
}

.message-wrapper.user .message-text :deep(tr:nth-child(even)) {
  background-color: rgba(255, 255, 255, 0.05);
}

/* 滚动条样式 */
.chat-messages::-webkit-scrollbar,
.history-list::-webkit-scrollbar {
  width: 6px;
}

.chat-messages::-webkit-scrollbar-track,
.history-list::-webkit-scrollbar-track {
  background: transparent;
}

.chat-messages::-webkit-scrollbar-thumb,
.history-list::-webkit-scrollbar-thumb {
  background: #d1d5db;
  border-radius: 3px;
}

.chat-messages::-webkit-scrollbar-thumb:hover,
.history-list::-webkit-scrollbar-thumb:hover {
  background: #9ca3af;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .sidebar {
    position: absolute;
    left: 0;
    top: 0;
    height: 100%;
    z-index: 1000;
    box-shadow: 4px 0 12px rgba(0, 0, 0, 0.15);
  }

  .sidebar-collapsed {
    transform: translateX(-100%);
  }

  .message-container {
    max-width: 90%;
  }

  .suggestions {
    grid-template-columns: 1fr;
  }

  .welcome-content h2 {
    font-size: 24px;
  }

  .welcome-content p {
    font-size: 16px;
  }
}
</style>