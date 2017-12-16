<template>
        <div class="container">
            <div class="container">
            <el-card class="box-card">
              <div slot="header" class="clearfix card-header">
                <span>Register</span>
              </div>
              <el-form ref="registration" status-icon :model="registration" :rules="rules">
                <el-form-item label="Username" prop="username">
                  <el-input v-model="registration.username" placeholder="Enter username"></el-input>
                </el-form-item>
                <el-form-item label="Email" prop="email">
                  <el-input v-model="registration.email" placeholder="Enter email"></el-input>
                </el-form-item>
                <el-form-item label="Password"  prop="password">
                  <el-input type="password" v-model="registration.password"  placeholder="Enter password" auto-complete="off"></el-input>
                </el-form-item>
                <el-form-item label="Confirm password"  prop="password2">
                  <el-input type="password" v-model="registration.password2"  placeholder="Enter the same password" auto-complete="off"></el-input>
                </el-form-item>

                <el-form-item style="text-align: right;">
                  <el-button type="danger" plain @click="onReset('registration')">Reset</el-button>
                  <el-button type="primary" @click="onSubmit('registration')"><font-awesome-icon icon="user-plus" /> Register</el-button>
                </el-form-item>
              </el-form>
            </el-card>
        </div>
      </div>
    </div>
</template>

<script>
import auth from '@/auth'
export default {
  name: 'register',
  data () {
    var validatePass2 = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('Please input the password again'))
      } else if (value !== this.registration.password) {
        callback(new Error('Two inputs don\'t match!'))
      } else {
        callback()
      }
    }
    return {
      registration: {
        username: '',
        email: '',
        password: '',
        password2: ''
      },
      rules: {
        username: [
          { required: true, message: 'Please input an username', trigger: 'blur' },
          { min: 3, message: 'Your username must be longer than 3 characters', trigger: 'blur' }
        ],
        email: [
          { required: true, message: 'Please input an email', trigger: 'blur' }
        ],
        password: [
          { required: true, message: 'Please input a password', trigger: 'blur' },
          { min: 6, message: 'Your password must be longer than 6 characters', trigger: 'blur' }
        ],
        password2: [
          { required: true, message: 'Please input the same password', trigger: 'blur' },
          { validator: validatePass2, trigger: 'blur' }
        ]
      },
      show: true
    }
  },
  methods: {
    onSubmit (formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.$http.post('http://localhost:1323/signup', {username: this.registration.username, password: this.registration.password, email: this.registration.email})
          .then((data) => {
            auth.login(this, {email: data.body.email, password: this.registration.password}, '/feed')
          })
        } else {
          console.log('error submit!!')
          return false
        }
      })
    },
    onReset (formName) {
      this.$refs[formName].resetFields()
    }
  }
}
</script>