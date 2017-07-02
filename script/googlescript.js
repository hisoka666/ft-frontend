function onSignIn(googleUser) {
  var id_token = googleUser.getAuthResponse().id_token;
  var url = "/login?idtoken=" + id_token;
  $.get(url)
  .done(function(data){
    if ($.trim(data) == "no-access"){
      signOut();
      alert("Maaf Anda tidak terdaftar sebagai staf. Silahkan hubungi admin");
     } else {
      // alert("Token adalah: " + $.trim(data));
      localStorage.setItem("token", $.trim(data));
      $("#signoutbut").show();
     }
  })
  .fail(function(err){
    alert("Error: " + err)
  });
};

// function getMain() {
//   var token = localStorage.getItem("token")
//   var url = "/getmain?token=" + token
//   $.get(url)
//   .done(function(data){

//   })

function signOut() {
   var auth2 = gapi.auth2.getAuthInstance();
   auth2.signOut().then(function () {
   console.log('User signed out.');
   localStorage.clear();
   $("#signoutbut").hide();
   });
};