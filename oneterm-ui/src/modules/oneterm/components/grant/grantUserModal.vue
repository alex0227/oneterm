<template>
  <a-modal
    :title="title"
    :visible="visible"
    destroyOnClose
    @ok="handleOk"
    @cancel="handleCancel"
  >
    <EmployeeTransfer
      :isDisabledAllCompany="true"
      v-if="type === 'depart'"
      uniqueKey="acl_rid"
      ref="employeeTransfer"
      :height="350"
    />
    <RoleTransfer
      v-if="type === 'role'"
      app_id="oneterm"
      :height="350"
      ref="roleTransfer"
    />
  </a-modal>
</template>

<script>
import EmployeeTransfer from '@/components/EmployeeTransfer'
import RoleTransfer from '@/components/RoleTransfer'

export default {
  name: 'GrantUserModal',
  components: {
    EmployeeTransfer,
    RoleTransfer
  },
  props: {
    customTitle: {
      type: String,
      default: '',
    },
  },
  data() {
    return {
      visible: false,
      type: 'depart',
    }
  },
  computed: {
    title() {
      if (this.type === 'depart') {
        return this.$t('oneterm.assetList.grantUserOrDep')
      }
      return this.$t('oneterm.assetList.grantRole')
    },
  },
  methods: {
    open(type) {
      this.visible = true
      this.type = type
    },
    handleOk() {
      let params
      if (this.type === 'depart') {
        params = this.$refs.employeeTransfer.getValues()
      }
      if (this.type === 'role') {
        params = this.$refs.roleTransfer.getValues()
      }
      this.handleCancel()
      this.$emit('handleOk', params, this.type)
    },
    handleCancel() {
      this.visible = false
    },
  },
}
</script>
