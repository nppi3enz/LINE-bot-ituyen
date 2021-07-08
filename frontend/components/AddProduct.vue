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
        <Button type="primary">
          Submit
        </Button>
      </Form>
    </div>
  </div>
</template>

<script>
import { StreamBarcodeReader } from 'vue-barcode-reader'
import { Form, Button, Input, FormItem, InputNumber, DatePicker } from 'element-ui'

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
    }
  }
}
</script>
