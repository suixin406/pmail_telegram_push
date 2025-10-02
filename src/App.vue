<template>
  <div class="pmail-telegram-push-settings">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>Telegram 推送设置</span>
          <div style="float: right">
            <el-badge :value="saved ? 1 : ''" class="item">
              <el-button
                type="primary"
                :disabled="!botInfo?.bot_link"
                @click="contactBot"
                style="float: right"
              >
                <i class="iconfont icon-telegram"> 联系机器人</i>
              </el-button>
            </el-badge>
          </div>
        </div>
      </template>

      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="180px"
        label-position="left"
        v-loading="loading"
      >
        <el-form-item label="Telegram 聊天 ID" prop="chat_id">
          <el-input
            v-model="formData.chat_id"
            placeholder="请输入 Telegram Chat ID，置空则禁用推送"
          />
        </el-form-item>

        <el-form-item label="显示邮件内容">
          <el-switch v-model="formData.show_content" />
        </el-form-item>

        <el-form-item label="显示邮件内容时添加防剧透">
          <el-switch v-model="formData.spoiler_content" />
        </el-form-item>

        <el-form-item label="发送附件">
          <el-switch v-model="formData.send_attachments" />
        </el-form-item>

        <el-form-item label="禁用链接预览">
          <el-switch v-model="formData.disable_link_preview" />
        </el-form-item>
      </el-form>
      <el-form-item>
        <el-button type="primary" @click="confirmSubmit" :loading="loading" style="margin: 0 auto">
          <i class="iconfont icon-save-line"> 保存设置</i>
        </el-button>
      </el-form-item>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import axios from 'axios'
import './resize.js'

const formRef = ref(null)
const loading = ref(false)
const formData = ref({
  chat_id: '',
  show_content: true,
  spoiler_content: true,
  send_attachments: true,
  disable_link_preview: true,
})

const rules = {
  chat_id: [
    {
      required: false,
      message: '请输入 Chat ID， 置空则禁用推送',
      trigger: 'blur',
      whitespace: false,
    },
  ],
}

const botInfo = ref({
  username: '',
  first_name: '',
  bot_link: '',
})

// 获取设置数据
const fetchSettings = async () => {
  try {
    loading.value = true
    const response = await axios.get('/api/plugin/settings/pmail_telegram_push/setting')
    if (response.data && response.data.code === 0) {
      formData.value = {
        ...formData.value,
        ...response.data.data,
      }
    } else {
      ElMessage.error(response.data.message || '获取设置失败')
    }
  } catch (error) {
    console.error('Failed to fetch settings:', error)
    ElMessage.error('获取设置失败')
  } finally {
    loading.value = false
  }
}

const saved = ref(false)

// 提交表单
const submitForm = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    const response = await axios.post(
      '/api/plugin/settings/pmail_telegram_push/submit',
      formData.value,
    )
    if (response.data && response.data.code === 0) {
      ElMessage.success('设置保存成功')
      if (formData.value.chat_id) {
        saved.value = true
      }
    } else {
      ElMessage.error(response.data.message || '保存设置失败')
    }
  } catch (error) {
    console.error('Failed to save settings:', error)
    if (error.response && error.response.data.message) {
      ElMessage.error(error.response.data.message)
    } else {
      ElMessage.error('保存设置失败')
    }
  } finally {
    loading.value = false
  }
}

// 确认提交表单
const confirmSubmit = () => {
  ElMessageBox.confirm('确认保存设置吗？', '确认', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(() => {
    submitForm()
  })
}

const getBotInfo = async () => {
  try {
    const response = await axios.get('/api/plugin/settings/pmail_telegram_push/bot')
    if (response.data && response.data.code === 0) {
      botInfo.value = {
        ...botInfo.value,
        ...response.data.data,
      }
    } else {
      ElMessage.error(response.data.message || '获取机器人信息失败')
    }
  } catch (error) {
    console.error('Failed to get bot:', error)
    ElMessage.error('获取机器人信息失败')
  }
}

// 联系机器人
const contactBot = () => {
  if (botInfo.value.bot_link) {
    saved.value = false
    window.open(botInfo.value.bot_link, '_blank')
  } else {
    ElMessage.error('机器人链接不存在')
  }
}

onMounted(() => {
  fetchSettings()
  getBotInfo()
})
</script>

<style scoped>
.pmail-telegram-push-settings {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.box-card {
  margin-top: 20px;
}

.card-header {
  font-size: 18px;
  font-weight: bold;
}
</style>
