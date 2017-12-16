<template>
        <div class="container">
            <div class="container">
            <el-card class="box-card">
              <div slot="header" class="clearfix card-header">
                <span>Login</span>
              </div>
              <el-alert type="error" title="ERROR!" description="Your credentials are not recognized by our system." show-icon v-show="showError"></el-alert>
              <el-form ref="login" :model="login" :rules="rules">
                <el-form-item label="Email" prop="email">
                  <el-input v-model="login.email" placeholder="Enter email"></el-input>
                </el-form-item>
                <el-form-item label="Password"  prop="password">
                  <el-input type="password" v-model="login.password"  placeholder="Enter password" auto-complete="off"></el-input>
                </el-form-item>

                <el-form-item style="text-align: right;">
                  <el-button type="danger" plain @click="onReset('login')">Reset</el-button>
                  <el-button type="primary" @click="onSubmit('login')"><font-awesome-icon icon="sign-in-alt" /> Login</el-button>
                </el-form-item>
              </el-form>
            </el-card>
        </div>
      </div>
    </div>
</template>

<script>
import auth from '../../auth'
import FontAwesomeIcon from '@fortawesome/vue-fontawesome'

export default {

  name: 'login',
  data () {
    return {
      showError: false,
      login: {
        password: '',
        email: ''
      },
      rules: {
        email: [
          { type: 'email', required: true, message: 'Please input an email', trigger: 'blur' }
        ],
        password: [
          { required: true, message: 'Please input a password', trigger: 'blur' },
          { min: 6, message: 'Password must be at least 6 characters long !', trigger: 'blur' }
        ]
      },
      show: true
    }
  },
  methods: {
    onSubmit (formName) {
      this.showError = false
      this.$refs[formName].validate((valid) => {
        if (valid) {
          auth.login(this, {email: this.login.email, password: this.login.password}, '/feed')
        } else {
          console.log('error submit!!')
          return false
        }
      })
    },
    onReset (formName) {
      this.$refs[formName].resetFields()
      this.showError = false
    }
  },
  components: {
    FontAwesomeIcon
  }
}
</script>