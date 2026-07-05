import axiosBase from "axios"

const axios = axiosBase.create({
  headers: {
    'Content-Type': 'application/json',
    'X-Requested-With': 'XMLHttpRequest'
  },
  responseType: 'json'  
});

class API {

  static post(url,args) {
    return axios.post(url,args)
  }

  static patch(url,args) {
    return axios.patch(url,args)
  }

  static delete(url,args) {
    return axios.delete(url,args)
  }
}

export default API;