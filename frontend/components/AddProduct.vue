<template>
    <div>
    <h1>Add Product</h1>
    <div v-if="!scanForm" id="scan">
        Please scan barcode before add product.
    <StreamBarcodeReader
      @decode="(a, b, c) => onDecode(a, b, c)"
      @loaded="() => onLoaded()"
    ></StreamBarcodeReader>
    Input Value: {{ barcode || "Nothing" }}
    </div>
    <div v-if="scanForm" id="input">
      <Form ref="form" :model="form" label-width="120px">
        <FormItem label="Barcode">
          <Input placeholder="Barcode" v-model="barcode" :disabled="true" />
        </FormItem>
        <FormItem label="Product name">
          <Input placeholder="Product Name" v-model="productName" />
        </FormItem>
        <FormItem label="Product List">
          <!-- <Input placeholder="Product List" v-model="text" /> -->
          <el-select v-model="productList" placeholder="Select">
            <el-option
            v-for="item in options"
            :key="item.value"
            :label="item.label"
            :value="item.value">
            </el-option>
          </el-select>
        </FormItem>
        <FormItem label="Expire Date">
          <div v-if="hasExpired">
            <DatePicker
                placeholder="Expire Date"
                type="date"
                :picker-options="pickerOptions"
                v-model="expireDate"
            />
          </div>
          <el-switch
            v-model="hasExpired"
            active-text="ON"
            inactive-text="OFF">
          </el-switch>
        </FormItem>
        <FormItem label="Quantitys">
          <InputNumber
            v-model="quantity"
            :min="1"
          />
        </FormItem>
        <Button type="primary" @click="submit()">
          Submit
        </Button>
      </Form>
    </div>
  </div>
</template>

<script>
import { StreamBarcodeReader } from 'vue-barcode-reader'
import { Form, Button, Input, FormItem, InputNumber, DatePicker } from 'element-ui'
import liff from '@line/liff'

liff.init({ liffId: '1656205141-4MGkXWrz' })
export default {
//   name: "HelloWorld",
  components: {
    StreamBarcodeReader,
    Input,
    Form,
    FormItem,
    Button,
    InputNumber,
    DatePicker
    // Row,
    // Col
  },
  data () {
    return {
    //   text: '',
      scanForm: false,
      id: null,
      barcode: '',
      productName: '',
      productList: '1',
      hasExpired: true,
      expireDate: '',
      quantity: 0,
      pickerOptions: {
        disabledDate (time) {
          return time.getTime() < Date.now()
        }
        // shortcuts: [{
        //   text: 'Today',
        //   onClick (picker) {
        //     picker.$emit('pick', new Date());
        //   }
        // }]
      },
      options: [{
        value: '1',
        label: 'Fridge'
      }]
    }
  },
  props: {
    // msg: String,
  },
  created () {
    this.$axios.get('https://staging-ituyen.herokuapp.com', {
      headers: {
        'Access-Control-Allow-Origin': '*',
        'Content-Type': 'application/Json'
      }
    }).then(response => console.log(response))
    this.getNow()
  },
  methods: {
    onDecode (a, b, c) {
      console.log(a, b, c)
      this.barcode = a
      this.scanForm = true
      if (this.id) {
        clearTimeout(this.id)
      }
      this.id = setTimeout(() => {
        if (this.text === a) {
          this.text = ''
        }
      }, 5000)
    },
    onLoaded () {
      console.log('load')
    },
    getNow () {
      const today = new Date()
      const date = today.getFullYear() + '-' + (today.getMonth() + 1) + '-' + today.getDate()
      this.expireDate = date
    },
    async submit () {
      // alert('submit')
      const data = await this.$axios.$post('https://staging-ituyen.herokuapp.com/product/create', {
        name: this.productName,
        barcode: this.barcode,
        expire_date: this.expireDate,
        quantity: this.quantity
      }).then((response) => {
        console.log(data.response.status)
        liff.closeWindow()
      }).catch(function (error) {
        if (error.response) {
          console.log(error.response.data)
          console.log(error.response.status)
          console.log(error.response.headers)
          alert(`Error: ${error.response.data.message}`)
        }
      })
      // return {posts: data}
    }

  },
  render: {
    static: {
      setHeaders (res) {
        res.setHeader('X-Frame-Options', 'ALLOWALL')
        res.setHeader('Access-Control-Allow-Origin', '*')
        res.setHeader('Access-Control-Allow-Methods', 'GET')
        res.setHeader('Access-Control-Allow-Headers', 'Origin, X-Requested-With, Content-Type, Accept')
      }
    }
  }
}
</script>
