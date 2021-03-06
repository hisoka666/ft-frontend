function onSignIn(googleUser) {
  var id_token = googleUser.getAuthResponse().id_token;
  var url = "/login?idtoken=" + id_token;
  $("#loading-animation").modal({backdrop: 'static', keyboard: false})
  $.get(url)
  .done(function(data){
    var jos = JSON.parse(data)
    if (jos.token == ""){
      signOut();
      alert("Maaf Anda tidak terdaftar sebagai staf. Silahkan hubungi admin");
      $("#signinbut").show();
     } else {
      localStorage.setItem("token", jos.token);
      localStorage.setItem("user", jos.modal);
      console.log("user adalah: "+ localStorage.getItem("user"))
      $("#dokter").html(jos.modal);
      $("#navbar").html(jos.script);
     }
     $("#loading-animation").modal('hide')
  })
  .fail(function(err){
    alert("Error: " + err)
  });

  var profile = googleUser.getBasicProfile();
  // $("#dokter").html(profile.getName());
  $("#email").attr({"value" : profile.getEmail()});
  $("#welcome").show();
  $("#signoutbut").show();
  $("#signinbut").hide();
};

function signOut() {
   var auth2 = gapi.auth2.getAuthInstance();
   auth2.signOut().then(function () {
   console.log('User signed out.');
   localStorage.clear();
   $("#welcome").hide();
   $("#signoutbut").hide();
   $("#navbar").html("");
   $("#signinbut").show();
   });
};
