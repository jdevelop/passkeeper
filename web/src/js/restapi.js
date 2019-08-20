import axios from "axios";

export const REST = {
  UpdateCardPassword(cardPassword, ok, fail) {
    axios
      .put(`${process.env.VUE_APP_API_URL}/cardpassword`, cardPassword, {})
      .then(data => ok(data))
      .catch(e => fail(e));
  },
  GeneratePasswords(done) {
    axios
      .get(`${process.env.VUE_APP_API_URL}/generate`, {}, {})
      .then(data => done(data));
  },
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
  },
  Backup() {
    axios
      .get(`${process.env.VUE_APP_API_URL}/backup`, { responseType: "blob" })
      .then(response => {
        const url = window.URL.createObjectURL(new Blob([response.data]));
        const link = document.createElement("a");
        link.href = url;
        link.setAttribute("download", "passkeeper-credentials.json");
        document.body.appendChild(link);
        link.click();
      });
  }
};
