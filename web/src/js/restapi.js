import axios from "axios";

export const REST = {
  SaveCredentials(credentials, done) {
    axios
      .put(`${process.env.VUE_APP_API_URL}/add`, credentials, {})
      .then(data => done(data));
  },
  ListCredentials(done) {
    axios
      .get(`${process.env.VUE_APP_API_URL}/list`, {}, {})
      .then(data => done(data.data));
  },
  RemoveCredentials(id, done) {
    axios
      .delete(`${process.env.VUE_APP_API_URL}/${id}`, {}, {})
      .then(data => done(data.data));
  }
};
