<template>
  <div>
    <header class="bg-white shadow">
      <div class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
        <h1 class="text-3xl font-bold text-gray-900">
          Remove Product
        </h1>
      </div>
    </header>
    <main>
      <div class="max-w-7xl mx-auto py-0 sm:px-6 lg:px-8">
        <!-- Replace with your content -->
        <div class="px-4 py-6 sm:px-0">
          <div v-if="step == 1" class="border-4 border-gray-200 rounded-lg p-5">
            <div class="loader ease-linear rounded-full border-8 border-t-8 border-gray-200 h-32 w-32"></div>
          </div>
          <div v-if="step == 2" class="border-4 border-gray-200 rounded-lg p-2">
            <StreamBarcodeReader
              @decode="(a, b, c) => onDecode(a, b, c)"
              @loaded="() => onLoaded()"
            />
          </div>
          <div v-if="step == 3" class="border-4 border-gray-200 rounded-lg p-2">
            <div class="px-2 py-2 bg-white sm:p-6">
              <div class="grid gap-4">
                <div class="col-span-6 sm:col-span-3">
                  <label for="barcode" class="block text-sm font-medium text-gray-700">Barcode</label>
                  <input
                    id="barcode"
                    v-model="barcode"
                    type="text"
                    name="barcode"
                    class="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md disabled:opacity-50"
                    :disabled="true"
                  >
                </div>
                <div class="col-span-6 sm:col-span-3">
                  <label for="productName" class="block text-sm font-medium text-gray-700">Product Name</label>
                  <input
                    id="productName"
                    v-model="productName"
                    type="text"
                    name="productName"
                    class="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md disabled:opacity-50"
                    :disabled="true"
                  >
                </div>
                <div class="col-span-6 sm:col-span-3">
                  <label for="quantity" class="block text-sm font-medium text-gray-700">Quantity</label>
                  <div class="custom-number-input">
                      <div class="flex flex-row h-10 w-full rounded-lg relative bg-transparent mt-1">
                        <button
                          data-action="decrement"
                          class="bg-gray-300 text-gray-600 hover:text-gray-700 hover:bg-gray-400 h-full w-20 rounded-l cursor-pointer outline-none"
                          @click="increment(-1)"
                        >
                          <span class="m-auto text-2xl font-thin">−</span>
                        </button>
                        <input
                          type="number"
                          class="outline-none focus:outline-none text-center w-full font-semibold text-md hover:text-black focus:text-black  md:text-basecursor-default flex items-center outline-none"
                          name="quantity"
                          v-model="quantity"
                        />
                        <button
                          data-action="increment"
                          class="bg-gray-300 text-gray-600 hover:text-gray-700 hover:bg-gray-400 h-full w-20 rounded-r cursor-pointer"
                          @click="increment(1)"
                        >
                          <span class="m-auto text-2xl font-thin">+</span>
                        </button>
                      </div>
                  </div>
                </div>
              </div>
            </div>
            <div v-if="step == 3" class="px-4 py-3 bg-gray-50 text-right sm:px-6">
              <button @click="submitExpiry()" type="submit" class="inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-emerald-600 hover:bg-emerald-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-emerald-500">
                Submit
              </button>
              <button @click="reset()" type="reset" class="inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-black bg-gray-300 hover:bg-gray-400 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-300">
                Reset
              </button>
            </div>
          </div>
          <div v-if="step == 9" class="border-4 border-gray-200 rounded-lg p-2 text-center">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-20 w-20 fill-current text-emerald-600 icon" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
            </svg>
            Add Complete
          </div>
          <div v-if="step == -1" class="border-4 border-gray-200 rounded-lg p-2 text-center">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-20 w-20 fill-current text-red-600 icon" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
            </svg>
            เกิดข้อผิดพลาด : {{ errorMsg }}
          </div>
        </div>
        <!-- /End replace -->
      </div>
    </main>
  </div>
</template>

<script>
import { StreamBarcodeReader } from 'vue-barcode-reader'
// import liff from '@line/liff'

export default {
  components: {
    StreamBarcodeReader
  },
  data () {
    return {
      step: 1,
      userID: null,
      barcode: '',
      productName: '',
      productList: '1',
      hasExpired: true,
      expireDate: '',
      quantity: 1,
      pickerOptions: {
        disabledDate (time) {
          return time.getTime() < Date.now()
        }
      },
      options: [{
        value: '1',
        label: 'Fridge'
      }],
      minDate: '',
      maxDate: '',
      accessToken: null,
      browserLanguage: null,
      sdkVersion: null,
      isInClient: null,
      isLoggedIn: null,
      deviceOS: null,
      errorMsg: ''
    }
  },
  async mounted () {
    await liff.init({ liffId: '1656205141-1QNAezQL' })
      .then(() => {
        if (!liff.isLoggedIn()) {
        //   liff.login()
        } else {
          this.step++
        }
      })
    this.getNow()
  },
  methods: {
    onDecode (a, b, c) {
      // console.log(a, b, c)
      this.barcode = a
      this.checkBarcode(this.barcode)
      this.step = 1
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
      const date = today.getFullYear() + '-' + String(today.getMonth() + 1).padStart(2, '0') + '-' + String(today.getDate() + 1).padStart(2, '0')
      this.expireDate = date
      this.minDate = date
    },
    increment (n) {
      if (this.quantity + n > 0) {
        this.quantity += n
      }
    },
    async checkBarcode (barcode) {
      await this.$axios.$get(`/api/expiry?barcode=${barcode}`).then((response) => {
        // console.log(response.data)
        if (response.data != null) {
          const firstItem = response.data[0]
          this.barcode = firstItem.product.barcode
          this.productName = firstItem.product.name
          if (firstItem.quantity > 1) {
            this.step = 3
          } else {
            this.deleteExpiry()
          }
        } else {
          this.errorMsg = 'ไม่ได้สามารถเชื่อมต่อกับฐานข้อมูลได้'
          this.step = -1
        }
      })
    },
    // async submit () {
    //   this.step = 1
    //   await this.$axios.$post('/api/product/create', {
    //     name: this.productName,
    //     barcode: this.barcode,
    //     expire_date: this.expireDate,
    //     quantity: this.quantity
    //   }).then((response) => {
    //     this.step = 9
    //     liff.closeWindow()
    //   }).catch(function (error) {
    //     if (error.response) {
    //       console.log(error.response.data)
    //       console.log(error.response.status)
    //       console.log(error.response.headers)
    //       this.errorMsg = error.response.data.message
    //       this.step = -1
    //     }
    //   })
    // },
    async deleteExpiry () {
      this.step = 1
      self = this
      await this.$axios.request('/api/expiry', {
        data: {
          barcode: self.barcode,
          quantity: self.quantity
        },
        method: 'delete'
      }).then((response) => {
        liff.sendMessages([
          {
            type: 'text',
            text: 'เช็ควันหมดอายุ'
          }
        ])
        self.step = 9
        liff.closeWindow()
      }).catch(function (error) {
        if (error.response) {
          self.errorMsg = error.response.data.message
          self.step = -1
        }
      })
    },
    reset () {
      this.productName = ''
      this.productList = '1'
      this.hasExpired = true
      this.expireDate = ''
      this.quantity = 1
    },
    login () {
      liff.login({ redirectUri: 'https://ituyen.herokuapp.com/remove-product' })
    }
  }
}
</script>
<style scoped>
/* Toggle A */
input:checked ~ .dot {
  transform: translateX(100%);
  background-color: #48bb78;
}

/* Toggle B */
input:checked ~ .dot {
  transform: translateX(100%);
  background-color: #48bb78;
}
input[type='number']::-webkit-inner-spin-button,
input[type='number']::-webkit-outer-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

.custom-number-input input:focus {
  outline: none !important;
}

.custom-number-input button:focus {
  outline: none !important;
}
.icon{
  display: block;
  margin: auto;
}
.loader {
  border-top-color: #3498db;
  -webkit-animation: spinner 1.5s linear infinite;
  animation: spinner 1.5s linear infinite;
  margin: 0 auto;
}

@-webkit-keyframes spinner {
  0% { -webkit-transform: rotate(0deg); }
  100% { -webkit-transform: rotate(360deg); }
}

@keyframes spinner {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>
