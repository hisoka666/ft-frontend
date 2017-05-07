function onSignIn(googleUser) {
  var id_token = googleUser.getAuthResponse().id_token;
  $.ajax({
    url: 'http://localhost:8080/test',
    type: 'post',
    data: {idtoken : id_token},
    success : function(data) {
      $('#ubah').html(data);

    },
  });
};

function signOut() {
   var auth2 = gapi.auth2.getAuthInstance();
   auth2.signOut().then(function () {
   console.log('User signed out.');
   });
};
