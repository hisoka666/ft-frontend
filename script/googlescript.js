function onSignIn(googleUser) {
  var id_token = googleUser.getAuthResponse().id_token;
  var url = "/login?idtoken=" + id_token;
  $.get(url)
  .done(function(data){

    var jos = JSON.parse(data)

    if (jos.token == ""){
      signOut();
      alert("Maaf Anda tidak terdaftar sebagai staf. Silahkan hubungi admin");
      $("#signinbut").show();
     } else {

      localStorage.setItem("token", jos.token);
      $("#navbar").html(jos.script)
      $("#signoutbut").show();
      $("#signinbut").hide();
     }
  })
  .fail(function(err){
    alert("Error: " + err)
  });
};

function signOut() {
   var auth2 = gapi.auth2.getAuthInstance();
   auth2.signOut().then(function () {
   console.log('User signed out.');
   localStorage.clear();
   $("#signoutbut").hide();
   $("#navbar").html("");
   $("#signinbut").show();
   });
};