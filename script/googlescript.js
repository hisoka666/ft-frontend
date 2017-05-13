function onSignIn(googleUser) {
  var id_token = googleUser.getAuthResponse().id_token;
  $.ajax({
    url: 'http://2.igdsanglah.appspot.com/test',
    type: 'post',
    data: {idtoken : id_token},
    success : function(data) {
      $('#ubah').html(data);
      localStorage.setItem("token", data)

    },
  });
};

function signOut() {
   var auth2 = gapi.auth2.getAuthInstance();
   auth2.signOut().then(function () {
   console.log('User signed out.');
   });
};
