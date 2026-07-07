import axiosBase from "axios"

const axios = axiosBase.create({
  headers: {
    'Content-Type': 'application/json',
    'X-Requested-With': 'XMLHttpRequest'
  },
  responseType: 'json'
});

// The server reports failures as 4xx/5xx with a {"error": "..."} body.
// Log the detail here, then keep the promise rejected so callers' .then()
// (which assumes success) never runs on a failed request.
axios.interceptors.response.use(
  (res) => res,
  (err) => {
    const detail = err.response && err.response.data && err.response.data.error;
    console.error("API error:", detail || err.message);
    return Promise.reject(err);
  }
);

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