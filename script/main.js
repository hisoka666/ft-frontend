
$(document).ready(function(){
	var nocm = $("#nocm").val()
	
	$("#navbar").on("keyup", "#nocm", function(){
		// var str = $("#nocm").val();
		var value = $("#nocm").val().replace(/\s+/g, '');
		var reg = /^\d+$/
		// console.log(reg.test(value))
		if (reg.test(value) == false){
			$("#datapasien").html("Harap masukkan angka!")
			value = "";
		}else if (value == ""){
			$("#datapasien").html("Masukkan No. CM");
			value = "";
		} else if (value.length < 8){
			$("#datapasien").html("No. CM tidak lengkap");
			value = "";
		} else {
			$("#nocm").prop("disabled", true);
			nocm = value;
            token = localStorage.getItem("token");

			$.post("getcm", {
				nocm: nocm,
				token: token
			},
			function(data){
				var js = JSON.parse(data)

				if (js.token == "OK"){
					$("#datapasien").html(js.script);
					$("#nocm").prop("disabled", false);
				} else {
					alert("Terjadi kesalahan: " + js.token)
					$("#nocm").prop("disabled", false);
				}			    

			});
        }

	});
	
	$("#navbar").on("click", "#btnsub", function(event){		
		event.preventDefault();
		
		var nocm = $("#nocm").val().replace(/\s+/g, '');
		var namapts = $("#namapts").val();
		var diag = $("#diag").val();
		var ats = $("input[name='ats']:checked").val();
		var iki = $("input[name='iki']:checked").val();
		var shift = $("input[name='shift']:checked").val();
		var bagian = $("input[name='bagian']:checked").val();
		var baru = $("#baru").val();
		var dok = $("#email").val();

		switch (true) {
			case nocm === "":
			break;
			case namapts === "":
			break;
			case diag === "":
			$("#alertmsg").html("<div class=\"alert alert-danger alert-dismissable\"\>" +
		                            "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a\>" +
									"Diagnosis belum diisi" +
		                        "</div>");
			// alert("Diagnosis harus diisi");
			break;
			case ats == null :
			$("#alertmsg").html("<div class=\"alert alert-danger alert-dismissable\"\>" +
		                            "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a\>" +
									"ATS belum diisi" +
		                        "</div>");
			// alert("ATS harus diisi");
			break;
			case iki == null :
			$("#alertmsg").html("<div class=\"alert alert-danger alert-dismissable\"\>" +
		                            "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a\>" +
									"Golongan IKI belum diisi" +
		                        "</div>");
			// alert("Golongan IKI harus diisi");
			break;
			case shift == null :
			$("#alertmsg").html("<div class=\"alert alert-danger alert-dismissable\"\>" +
		                            "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a\>" +
									"Shift belum diisi" +
		                        "</div>");
			// alert("Shift harus diisi");
			break;
			case bagian == null:
			$("#alertmsg").html("<div class=\"alert alert-danger alert-dismissable\"\>" +
		                            "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a\>" +
									"Bagian belum diisi" +
		                        "</div>");
			// alert("Bagian harus diisi");
			break;
			default:
			$.post("inputdata",{
				token: localStorage.getItem("token"),
				nocm: nocm,
				namapts: namapts,
				diag: diag,
				ats: ats,
				iki: iki,
				shift: shift,
				bagian: bagian,
				dok: dok,
				baru: baru
			},
			function(data){
				var js = JSON.parse(data)
				if (js.token != "OK"){
					alert(js.script);
				}else{
					$("tbody.dafpts").prepend(js.script);
					$("#nocm").val('');
					$("#datapasien").html('');
					$("tbody tr:eq(100)").remove();
					refreshNumber();
					popModalWarning("Sukses", "Berhasil menambahkan data", "")
					
				}
			});
		}
	});

	$("#navbar").on("click", "#editbut", function(event){
		event.preventDefault();
		var link = $(this).offsetParent().children().first().html();
		var indexrow = $(this).closest("tr").index();
		token = localStorage.getItem("token");
			$.post("editentri",{
				token: token,
				link: link
			},
			function(data){
				var js = JSON.parse(data);
				$("#mymodal").html(js.script);
				var ats = 'input[name="modats"][value=' + js.content.ats + ']';
				var iki = 'input[name="modiki"][value=' + js.content.iki + ']';
				var shift = 'input[name="modshift"][value=' + js.content.shift + ']';
				var bagian = 'input[name="modbagian"][value=' + js.content.dept + ']';
			
				$('input[name="entri"]').val(js.content.link);
				$('input[name="namapasien"]').val(js.content.nama);
				$('input[name="diagnosis"]').val(js.content.diag);
				$('input[name="urutan"]').val(indexrow)
			
				if (js.content.bagian == ""){
 		    		$(ats).prop('checked', true);
   					$(iki).prop('checked', true);
   					$(shift).prop('checked', true);	
			
				}else{
 		    		$(ats).prop('checked', true);
   					$(iki).prop('checked', true);
   					$(shift).prop('checked', true);
   					$(bagian).prop('checked', true);

				};
				$("#mymodal").modal();
			});
	});

$('body').on('hide.bs.modal', "#mymodal" , function(){
	console.log("Removal mymodal works");
	$(this).removeData('bs.modal').find(".modal-content").empty();
	});
$('body').on('hide.bs.modal', "#modwarning" , function(){
	console.log("removal modwarning works");
	$(this).removeData('bs.modal');
	});

$("body").on("click", "#confirmedit", function(e){
		e.preventDefault();
		
		var link = $("#modentri").val();
		var namapts = $('input[name="namapasien"]').val();
		var diag = $('input[name="diagnosis"]').val();
		var ats = $("input[name='modats']:checked").val();
		var bagian = $("input[name='modbagian']:checked").val();
		var iki = $("input[name='modiki']:checked").val();
		var shift = $("input[name='modshift']:checked").val();
		var indexrow = $('input[name="urutan"]').val();
		var urutan = "tbody.dafpts tr:eq(" + $('input[name="urutan"]').val() + ")";

		
		switch (true) {
			case namapts === "":
			popModalWarning();
			break;
			case diag === "":
			popModalWarning();
			break;
			case ats == null :
			break;
			case iki == null :
			break;
			case shift == null :
			break;
			case bagian == null:
			break;
			default:
			$.post("confedit",{
				token: localStorage.getItem("token"),
				nocm: nocm,
				namapts: namapts,
				diag: diag,
				ats: ats,
				iki: iki,
				shift: shift,
				bagian: bagian,
				link: link
			},
			function(data){
				var js = JSON.parse(data)
				if (js.token != "OK"){
					popModalWarning("Kesalahan Pada Server", "Telah terjadi kesalahan pada server. Mohon ulangi proses sebelumnya");
				}else{
					$(urutan).replaceWith(js.script);
					console.log("Editing jalan")
					refreshNumber();
					$("#mymodal").modal('hide');
					popModalWarning("Edit Entri", "Berhasil mengubah entri", "");
					
				}
			});
		}

	});

var refreshNumber = function(){
	
$("tr").find(".nourut").each(function(index, elem){
		num = index + 1;
		$(elem).html(num)
	})
}

var popModalWarning = function(title, msg, prop){
		$(".modal-title").html(title);
		$("#message").html(msg);
		if (prop == ""){
			$("#extrabut").hide();
			$("#modwarning").modal();
		} else {
			$("#extrabut").html(prop);
			$("#extrabut").show();
			$("#modwarning").modal();
		}			
};

$("#navbar").on("click", "#delbut", function(e){
		e.preventDefault();
		var link = $(this).offsetParent().children().first().html();
		var indexrow = $(this).closest("tr").index();
		var urutan = "tbody.dafpts tr:eq(" + indexrow + ")";
		var token = localStorage.getItem("token");
		console.log("Link adalah: " + link)
		popModalWarning("Hapus Entri", "Yakin ingin menghapus entri ini?", "Hapus");
		$("body").one("click", "#extrabut", function(){
			
			$.post("confdel", {
				link: link,
				token: token
			}, function(data){
				var js = JSON.parse(data);
				if (js.token == "OK"){
					$(urutan).remove();
					console.log("Removed : "+ urutan);
					refreshNumber();
					$("#modwarning").modal('hide');
				}else{
					alert("Terjadi kesalahan");
				}
			});
		});
	});






$("#navbar").on("click", "#homebutton", function(e){
		token = localStorage.getItem("token");
		email = $("#email").val();
		e.preventDefault();
		$.post("/firstentries", {
			token: token,
			email: email
		},function(data){
			var js = JSON.parse(data);
			$("div#resep").hide();
			$("div#main").show();
			// console.log(js.script);
			if (js.token == "OK") {
				$("#inputnocm").show();
				$("#tabeliki").html("").hide();
				$(".diagram").hide();
				$("#ket-bulan").hide();
				$("#tabelutama").html(js.script);
				refreshNumber();
				// removeModal("#modwarning")
			}else{
				popModalWarning("Peringatan", "Terjadi kesalahan pada server. Hubungi admin", "");
			}
			

		})

	});

$("#navbar").on("click", "#tglbut", function(e){
	token = localStorage.getItem("token");
	var link = $(this).offsetParent().children().first().html();
	e.preventDefault();
	console.log(link);

	$.post("edittgl", {
		token: token,
		link: link
	}, function(data){
		var js = JSON.parse(data);
		if (js.token != "OK") {
			$("#mymodal").html(js.modal);
		}else{
			$("#mymodal").html(js.script);
			$("#inputdate").datepicker({
				dateFormat: "dd-mm-yy",
				altField: "#datesend",
				altFormat: "yy-mm-dd"
			})
		}
		$("#mymodal").modal();
	})
});

$("body").on("click", "#ubahtglbut", function(e){
	e.preventDefault();
	var token = localStorage.getItem("token");
	var tanggal = $("#datesend").val();
	var link = $("#linkubahtgl").val();
	console.log("tanggal adalah : " + tanggal);
	// var tglasli = $(this).find("#tglasli").html();
	$.post("confedittgl", {
		token: token,
		link: link,
		tanggal: tanggal
	}, function(data){
		var js = JSON.parse(data);
		// console.log("Isi dari script adalah: " + data.script)
		$("#mymodal").modal('hide');
		if (js.token != "OK") {			
			popModalWarning("Peringatan", "Terjadi kesalahan pada server. Hubungi admin","");
		} else {
			$("#tabelutama").html(js.script);
			popModalWarning("Sukses", "Berhasil mengubah tanggal","");
		}

	})
	// console.log(jQuery.type(tanggal));
	// $("#mymodal").modal('hide');
});

$("#navbar").on("click", "#makeresep", function(e){
	e.preventDefault();

	$.get("getptspage")
	.done(function(data){
		$("div#main").hide();
		var js = JSON.parse(data);
		$("div#resep").html(js.script).show();
		// $("#mymodal").html(js.script);
		// $("#mymodal").modal();
		$("#tgllahir").datepicker({
			dateFormat:"dd-mm-yy",
			changeMonth: true,
			changeYear: true,
			yearRange: "1900:2035",
			onSelect: function(value, ui){
				var today = new Date();
				// console.log("Tahun ini adalah: " + today.getFullYear());
				// console.log("Tahun yang dipilih adalah : " + ui.selectedYear)
				var umur = today.getFullYear() - ui.selectedYear;
				// console.log("Umur adalah: " + umur);
				$("#umur").val(umur);
			}
		});
	});

	$("body").on("click", "#rspnextbut", function(e){
		e.preventDefault();
		
		var nama = $("#namapts").val();
		var diag = $("#diag").val();
		var umur = $("#umur").val();
		var almt = $("#almt").val();
		var bb = $("#bb").val();
		var alergi = $("#alergi").val();
		var ptsid = $(".ptsid").val();
		console.log("Berat adalah: " + bb);
		if (bb == 0){
			$("#alertmsgobat").html("<div class=\"alert alert-danger alert-dismissable\"\>" +
			"<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a\>" +
			"Berat badan harus diisi!" +
			"</div>");
		}else{
			console.log("fired!")
		// $("#mymodal").modal();
		$.get("getprespage")
		.done(function(resep){
			var jso = JSON.parse(resep);
			$("div#resep").html(jso.script);
			$("#rspnamapts").html(nama);
			$("#rspdiag").html(diag);
			$("#rspumur").html(umur);
			$("#rspbb").html(bb);
			$("#rspalmt").html(almt);
			$("#rspalergi").html(alergi);
			$(".ptsid").html(ptsid);
			// $("#mymodal").modal();
		})
	}
		
	})
});

$("body").on("change", "input:radio[name=sediaan]:checked", function(){
	if (this.value == "1"){
		$("div.obat").hide();
		$("input.obat").val("");
		$("div.obat.tablet").show();
	}else if (this.value == "2"){
		// $("div.obat").not("sirup drop").hide();
		$("div.obat").hide();
		$("input.obat").val("");
		$("div.obat.sirup, div.obat.drop").show();
	}else{
		$("div.obat").hide();
		$("input.obat").val("");
		// $("div.obat").not("lainnya").hide();
		$("div.obat.lainnya").show();
	}
	$("div.rekom").show();
})
$("#navbar").on("click", "#inputobat", function(e){
	e.preventDefault();
	$.get("getinputobat")
	.done(function(data){
		var js = JSON.parse(data);
		if (js.token != "OK"){
			popModalWarning("Peringatan", "Gagal memperoleh template", "")
		}else{
			$("#mymodal2").html(js.script)
			$("#inputby").val($("#email").val());
			$("#mymodal2").modal();
		}

	})
})

$("body").on("click", ".btn.tablet.tambah", function(e){
	tambahElement("obat.tablet.sediaan", this)
});

$("body").on("click", ".btn.tablet.hapus", function(e){
	e.preventDefault();
	hapusElement("obat.tablet.sediaan", this);
});

$("body").on("click", ".btn.sirup.tambah", function(e){
	e.preventDefault();
	tambahElement("obat.sirup.sediaan", this);
});

$("body").on("click", ".btn.sirup.hapus", function(e){
	e.preventDefault();
	hapusElement("obat.sirup.sediaan",this)
});

$("body").on("click", ".btn.drop.tambah", function(e){
	e.preventDefault();
	tambahElement("obat.drop.sediaan", this);
});

$("body").on("click", ".btn.drop.hapus", function(e){
	e.preventDefault();
	hapusElement("obat.drop.sediaan",this)
});
$("body").on("click", ".btn.lainnya.tambah", function(e){
	e.preventDefault();
	tambahElement("obat.lainnya.sediaan", this);
});

$("body").on("click", ".btn.lainnya.hapus", function(e){
	e.preventDefault();
	hapusElement("obat.lainnya.sediaan",this)
});
var tambahElement = function(selector, elem){
	$("div."+ selector + ".col-sm-9").last().clone().appendTo(".form-group."+selector);
	$("div."+ selector + ".col-sm-3").last().clone().appendTo(".form-group."+selector);
	$("input.form-control."+selector).each(function(index, elem){
		// $(this).prop("disabled",true);
		$(elem).attr("name", index);
	});
	$("input.form-control."+selector).last().val("");
	$(elem).html("Hapus").removeClass("tambah").addClass("hapus");

}

var hapusElement = function(sel, elem){
	var index = $("."+sel+".col-sm-3").index($(elem).parent());
	$("."+sel+".col-sm-9").eq(index).remove();
	$("."+sel+".col-sm-3").eq(index).remove();
}


$("body").on("click", "#savdrug", function(e){
	e.preventDefault();
	data = {
		"merk": $("#mrkdgng").val(),
		"kand": $("#kand").val(),
		"mindose": $("#mindose").val(),
		"maxdose": $("#maxdose").val(),
		"tab" : convertSerialArray($("input.tablet").serializeArray()),
		"syr" : convertSerialArray($("input.sirup").serializeArray()),
		"drop": convertSerialArray($("input.drop").serializeArray()),
		"lainnya_sediaan": convertSerialArray($("input.lainnya.sediaan").serializeArray()),
		"lainnya": $("input.lainnya.bentuk").val(),
		"rekom" : $("#rekom").val(),
		"doc" : $("#email").val()
	}

	console.log("String json adalah : " + JSON.stringify(data))
	
	$.post("inputobat", {
		send: JSON.stringify(data),
		token: localStorage.getItem("token")
		}, function(data){
			var js = JSON.parse(data);
			if (js.token != "OK"){
				$("#mymodal").html(js.modal)
				$("#mymodal").modal()
			}else{
				popModalWarning("Sukses", "Berhasil menambahkan obat", "")
			}
		})
	$("#mymodal2").modal('hide');


});


convertSerialArray = function(arr){
	var r = [];
	for (i=0;i<arr.length;i++){
		r[i] = arr[i].value
	}
	return r
}
$("#navbar").on("click", "#bulanini", function(e){
	e.preventDefault();
	var now = new Date();
	var dateone = new Date(now.getFullYear(),now.getMonth(),1,8,0,0);
	var token = localStorage.getItem("token");
	$("div#resep").hide();
	if (now > dateone){
		$.post("getmonthly", {
			token: token,
			month: now.getMonth() + 1,
			year: now.getFullYear(),
			email: $("#email").val()
		}, function(data){
			var js = JSON.parse(data);
			$("div#main").show();
			pieChart(js.data, "")
			$("#inputnocm").hide();
			$(".diagram").show();
			$("#butpdfnow").show();
			// $("#tabeliki").html(js.modal).show();
			$("#tabeliki").html(js.modal)
			getSum();
			var ekstra = js.data.data2
			// console.log("p3k " + ekstra.p3k.length)
			if (jQuery.isEmptyObject(ekstra.p3k) == false){
				for (i=0;i<ekstra.p3k.length;i++){
					$("td.iki1").eq(parseInt(js.data.data2.p3k[i])-1).html("P3K").removeClass("jml")
				}
			}
			if (jQuery.isEmptyObject(ekstra.rapat) == false){
				for (i=0;i<ekstra.rapat.length;i++){
					$("td.iki1").eq(parseInt(js.data.data2.rapat[i])-1).html("Rapat").removeClass("jml")
				}
			}
			if (jQuery.isEmptyObject(ekstra.pelatihan) == false){
				for (i=0;i<ekstra.pelatihan.length;i++){
					$("td.iki1").eq(parseInt(js.data.data2.pelatihan[i])-1).html("Pelatihan").removeClass("jml")
				}
			}
			$("#tabeliki").show();
			$("#tabelutama").html(js.script);
		})
	}else{
		var blnlalu = new Date(now.getFullYear(), now.getMonth() - 1, 1)
		$.post("getmonthly", {
			token: token,
			month: blnlalu.getMonth() + 1,
			year: blnlalu.getFullYear(),
			email: $("#email").val()
		}, function(data){
			$("div#main").show();
			var js = JSON.parse(data);
			pieChart(js.data, "")
			$(".diagram").show();
			$("#inputnocm").hide();
			// $("#tabeliki").html(js.modal).show();
			$("#tabeliki").html(js.modal)
			getSum();
			var ekstra = js.data.data2
			// console.log("p3k " + ekstra.p3k.length)
			if (jQuery.isEmptyObject(ekstra.p3k) == false){
				for (i=0;i<ekstra.p3k.length;i++){
					$("td.iki1").eq(parseInt(js.data.data2.p3k[i])-1).html("P3K").removeClass("jml")
				}
			}
			if (jQuery.isEmptyObject(ekstra.rapat) == false){
				for (i=0;i<ekstra.rapat.length;i++){
					$("td.iki1").eq(parseInt(js.data.data2.rapat[i])-1).html("Rapat").removeClass("jml")
				}
			}
			if (jQuery.isEmptyObject(ekstra.pelatihan) == false){
				for (i=0;i<ekstra.pelatihan.length;i++){
					$("td.iki1").eq(parseInt(js.data.data2.pelatihan[i])-1).html("Pelatihan").removeClass("jml")
				}
			}
			$("#tabeliki").show();
			$("#tabelutama").html(js.script);
		})
	}


});

var getSum = function(){
	var sum = 0;
	var iki11 = 0;
	var iki12 = 0;
	for (i=0;i<16;i++){
		iki11 = iki11 + parseInt($("td.jml").eq(i).html())
	}
	for (i=16;i<32;i++){
		iki12 = iki12 + parseInt($("td.jml").eq(i).html())
	}
	var iki21 = iki11;
	var iki22 = iki12;
	for (i=32;i<47;i++){
		iki21 = iki21 + parseInt($("td.jml").eq(i).html())
	}
	for (i=47;i<62;i++){
		iki22 = iki22 + parseInt($("td.jml").eq(i).html())
	}
	
	var tot1 = iki21 * 0.0032;
	var tot2 = iki22 * 0.01;
	var totfinal = tot1 + tot2;
	$("#jmlpoin-1-1").html(iki11);
	$("#jmlpoin-1-2").html(iki12);
	$("#jmlpoin-2-1").html(iki21);
	$("#jmlpoin-2-2").html(iki22);
	$("#jmlxpoin1").html(tot1.toFixed(4));
	$("#jmlxpoin2").html(tot2.toFixed(2));
	$("#totalpoin").html("Keterangan: Total Poin untuk perhitungan IKI adalah: " + totfinal.toFixed(4));
	$("#ket-bulan").show();
}

$("#navbar").on("click", ".bcptgl", function(e){
	e.preventDefault();
	var token = localStorage.getItem("token");
	var tgl = $(this).html();
	console.log(tgl)
	$("div#resep").hide();
	$.post("getbcpmonth", {
		token: token,
		tgl: tgl,
		email: $("#email").val()
	}, function(data){
		var js = JSON.parse(data);
		$("div#main").show();
		// console.log("Jumlah interna adalah : " + JSON.stringify(js.data.data1))
		pieChart(js.data, tgl);
		$("#inputnocm").hide();
		$(".diagram").show();
		$("#butpdfnow").hide();
		// console.log("p3k adalah : " +)
		// $("#tabeliki").html(js.modal).show()
		$("#tabeliki").html(js.modal)
		getSum();
		var ekstra = js.data.data2
		// console.log("p3k " + ekstra.p3k.length)
		if (jQuery.isEmptyObject(ekstra.p3k) == false){
			for (i=0;i<ekstra.p3k.length;i++){
				$("td.iki1").eq(parseInt(js.data.data2.p3k[i])-1).html("P3K").removeClass("jml")
			}
		}
		if (jQuery.isEmptyObject(ekstra.rapat) == false){
			for (i=0;i<ekstra.rapat.length;i++){
				$("td.iki1").eq(parseInt(js.data.data2.rapat[i])-1).html("Rapat").removeClass("jml")
			}
		}
		if (jQuery.isEmptyObject(ekstra.pelatihan) == false){
			for (i=0;i<ekstra.pelatihan.length;i++){
				$("td.iki1").eq(parseInt(js.data.data2.pelatihan[i])-1).html("Pelatihan").removeClass("jml")
			}
		}
		$("#tabeliki").show();
		$("div.tabtitle").html("Tabel IKI " + tgl);
		$("#tabelutama").html(js.script);
	})
});

$("#navbar").on("click", ".createpdf", function(e){
	e.preventDefault();
	var token = localStorage.getItem("token");
	var tgl = $(this).html();
	var nama = "dr. " + $("#dokter").html();
	var email = $("#email").val();
	$("#emailpdf").val(email);
	$("#namapdf").val(nama);
	$("#tglpdf").val(tgl);
	$("#tokenpdf").val(token);
	$("#getpdf").prop("action", "/getpdf")
	$("#getpdf").submit();

});

$("#navbar").on("click", "#bul-ini-pdf", function(e){
	e.preventDefault();
	var token = localStorage.getItem("token");
	// var tgl = $(this).html();
	var nama = "dr. " + $("#dokter").html();
	var email = $("#email").val();
	var now = new Date();
	var dateone = new Date(now.getFullYear(),now.getMonth(),1,8,0,0);
	var tgl = ""
	if ( now > dateone){
		tgl = now.getFullYear().toString() + "/" + ("0" + (now.getMonth() + 1).toString()).slice(-2)
	}else{
		var blnlalu = new Date(now.getFullYear(), now.getMonth() - 1, 1)
		tgl = blnlalu.getFullYear().toString + "/" + ("0" + (blnlalu.getMonth() + 1).toString()).slice(-2)
	}
	$("#emailpdf").val(email);
	$("#namapdf").val(nama);
	$("#tglpdf").val(tgl);
	$("#tokenpdf").val(token);
	$("#getpdf").prop("action", "/getpdfnow")
	$("#getpdf").submit();
})
var pieChart = function(list, tgl){
	google.charts.load('current', {packages: ['corechart', 'bar']});
	google.charts.setOnLoadCallback(drawChart);
	google.charts.setOnLoadCallback(barChart);
	function drawChart(){
		var data = google.visualization.arrayToDataTable([
          ['Bagian', 'Jumlah Pasien'],
          ['Interna', list.data0.interna],
          ['Bedah', list.data0.bedah],
          ['Anak',  list.data0.anak],
          ['OBGYN', list.data0.obgyn],
          ['Saraf', list.data0.saraf],
          ['Anestesi', list.data0.anes],
          ['Psikiatri', list.data0.psik],
          ['THT', list.data0.tht],
          ['Kulit dan Kelamin', list.data0.kulit],
          ['Jantung', list.data0.jant],
          ['Umum', list.data0.umum],
          ['Mata', list.data0.mata],
          ['MOD', list.data0.mod],
        ]);

        var options = {
		  title: 'Jumlah Pasien PerBagian',
		  width: 800,
		  height: 400,
        };

        var chart = new google.visualization.PieChart(document.getElementById('piechart'));

        chart.draw(data, options);
	}

	function barChart(){
		var data = new google.visualization.DataTable();
	    data.addColumn('number', 'Tanggal Jaga');
		data.addColumn('number', 'Jumlah IKI 1');
		data.addColumn('number', 'Jumlah IKI 2');
		// var title = [['Tanggal Jaga', 'Jumlah IKI 1', 'Jumlah IKI 2']]
		// var js = JSON.parse(list.data1)
		var isi = []
		for (i=0;i<31;i++){
			// var isi = array[awal.tgl, awal.iki1, awal.iki2]
			var awal = list.data1[i]
			data.addRows([
				[awal.tgl, awal.iki1, awal.iki2]
			])
		}

		
		// var data = new google.visualization.arrayToDataTable(title);
        var options = {
			title: 'Jumlah Pasien Pertanggal Jaga',
			width: 800,
			hAxis : {
				title: "Tanggal Jaga",
			},
			vAxis: {
				title: "Jumlah Pasien"
			}
        };
		var chart = new google.visualization.ColumnChart(document.getElementById('barchart'));

        chart.draw(data, options);

	}
}

$('body').on('keyup', 'input.isianobat', function(){
	var ob = ""
	var ob = $(this).val()
	var token = localStorage.getItem("token")
	var ini = $(this).parents(".listobat").children('#listobat')
	// if (ob == "puyer") {
	// 	$.get("formpuyer")
	// 	.done(function(data){
	// 		js = JSON.parse(data)
	// 		ini.html(js.script)
	// 	})
	// }else{
	$.post("cariobt",{
		token: token,
		obat: ob
	}, function(data){
		js = JSON.parse(data)
		ini.html(js.script)
	})
	// }
})

$('body').on('click', 'a.addobatinfo', function(e){
	e.preventDefault();
	obatnam = $(this).find("#obatbaru").html()
	// console.log(obatnam)
	$.get("getinputobat")
	.done(function(data){
		js = JSON.parse(data)
		// console.log(js.script)
		$("#mymodal2").html(js.script)
		$("#inputby").val($("#email").val());
		$("#mymodal2").modal();
	})
});
$('body').on('click', 'a.getobatinfo', function(e){
	e.preventDefault();
	link = $(this).attr('href');
	bb = $("#rspbb").html();
	token = localStorage.getItem("token");
	ini = $(this).parent()
	inputini = $(this).parents('.listobat').children().find('input.isianobat')
	$.post("getobat",{
		token: token,
		link: link,
		berat: bb
	}, function(data){
		js = JSON.parse(data);
		console.log(js.modal)
		inputini.val(js.modal)
		// console.log(js.script);
		ini.html(js.script)
	})
})

$('body').on('click', 'button.tambahlistobat', function(e){
	e.preventDefault();
	$('div.template').clone().removeClass('template').addClass('listobat').prop('hidden', false).appendTo('div.form-group.main')
})

$('body').on('click', 'a#editobat', function(e){
	e.preventDefault();
	link = $(this).attr('href');
	console.log("this is link: " + link)
	$.post("getobatedit", {
		token: localStorage.getItem("token"),
		link: link
	},function(data){
		js = JSON.parse(data)
		// console.log("This is data: " + data)
		$("#mymodal2").html(js.script)
		$("#mrkdgng").val(js.data.merk)
		$("#kand").val(js.data.kand)
		$("#mymodal2").modal()
		// for (i = 0; i<js.data.syr.length; i++){
		// 	$("input.obat.sirup.sediaan").clone()
		// }
		// var lainnya = 'input[name="sediaan"][value="3"]'
		// var 
		if (js.data.tab[0] == "" && js.data.syr[0] == ""){
			// for (i=0;i<js.data.lainnya_sediaan.length; i++){
			// 	$('input.obat.lainnya.sediaan').html(js.data.lainnya_sediaan[i])
			// 	tambahElement("obat.lainnya.sediaan", this);
			// }
			$('input[name="sediaan"][value="3"]').prop('checked', true);
			$("div.obat").hide();
			$("input.obat").val("");
			$("div.obat.lainnya").show();
		}else if (js.data.syr[0] !== ""){
			// for (i=0;i<js.data.syr.length; i++){
			// 	$('input.obat.sirup')
			// }
			$('input[name="sediaan"][value="2"]').prop('checked', true)
			$("div.obat").hide();
			$("input.obat").val("");
			$("div.obat.sirup, div.obat.drop").show();
		}else {
			$('input[name="sediaan"][value="1"]').prop('checked', true)
			$("div.obat").hide();
			$("input.obat").val("");
			$("div.obat.tablet").show();
		}
		$("div.rekom").show();
		$("input#mindose").val(js.data.mindose)
		$("input#maxdose").val(js.data.maxdose)
		$("input#linkedit").val(link)
	})
})

$("body").on("click", "#savdrugedit", function(e){
	e.preventDefault();
	data = {
		"merk": $("#mrkdgng").val(),
		"kand": $("#kand").val(),
		"mindose": $("#mindose").val(),
		"maxdose": $("#maxdose").val(),
		"tab": $("input.tablet").val(),
		"syr": $("input.sirup").val(),
		"drop": $("input.drop").val(),
		"lainnya": $("input.lainnya.sediaan").val(),
		"tab" : convertSerialArray($("input.tablet").serializeArray()),
		"syr" : convertSerialArray($("input.sirup").serializeArray()),
		"drop": convertSerialArray($("input.drop").serializeArray()),
		"lainnya_sediaan": convertSerialArray($("input.lainnya.sediaan").serializeArray()),
		"lainnya": $("input.lainnya.bentuk").val(),
		"rekom" : $("#rekom").val(),
		"doc" : $("#email").val()
	}

	// console.log("String json adalah : " + JSON.stringify(data))
	// console.log($("#linkedit").val())
	$.post("inputobatedit", {
		send: JSON.stringify(data),
		token: localStorage.getItem("token"),
		link: $("input#linkedit").val()
		}, function(data){
			var js = JSON.parse(data);
			if (js.token != "OK"){
				$("#mymodal").html(js.modal)
				$("#mymodal").modal()
			}else{
				popModalWarning("Sukses", "Berhasil merubah data obat", "")
			}
		})
	$("#mymodal2").modal('hide');
});

$('body').on('keyup', 'input.obat-puyer', function(){
	var ob = ""
	var ob = $(this).val()
	var token = localStorage.getItem("token")
	var ini = $(this).parents(".list-puyer").children('#listobatpuyer')
	// ini.html(ob)
	$.post("cariobatpuyer",{
		token: token,
		obat: ob
	}, function(data){
		js = JSON.parse(data)
		ini.html(js.script)
		$("#obatbaru").html(js.modal)
	})
})

$('body').on('click', 'a.getpuyerinfo', function(e){
	e.preventDefault();
	link = $(this).attr('href');
	bb = $("#rspbb").html();
	token = localStorage.getItem("token");
	ini = $(this).parent()
	inputini = $(this).parents('.list-puyer').children().find('input.obat-puyer')
	$.post("getpuyer",{
		token: token,
		link: link,
		berat: bb
	}, function(data){
		js = JSON.parse(data);
		// console.log(inputini.val())
		inputini.val(js.modal)
		// console.log(js.data.dosis);
		ini.html("Rekomendasi: " + js.data.rekom + ", Dosis: " + js.data.dosis + ". <a id='editobat' href=" + js.data.link + ">Edit</a>").addClass("help-block")
	})
})

$('body').on('click', 'button.add-obat-puyer', function(e){
	e.preventDefault();
	$('div.template-puyer').clone().removeClass('template-puyer').prop('hidden', false).appendTo('div.form-group.main-puyer')
})
$('body').on('click', 'button.but-puyer', function(e){
	e.preventDefault()
	console.log("button pressed")
	var ini = $(this).parents(".listobat").children("#listobat")
	var iniinput = $(this).parents('div.input-group').children('input.isianobat')
	$.get("formpuyer")
		.done(function(data){
			js = JSON.parse(data)
			iniinput.remove()
			ini.html(js.script)
		})
})

$('body').on('click', 'button.del-obat-line', function(e){
	e.preventDefault();
	$(this).parents('.listobat').remove();
})

$('body').on('click', 'button#resepbut', function(e){
	e.preventDefault();
	var tgl = new Date()
	var bul = ("0" + (tgl.getMonth()+1).toString()).slice(-2)
	var thn = tgl.getFullYear().toString()
	var hari = ("0" + (tgl.getDate()+1).toString()).slice(-2)
 	var strDate = hari + "/" + bul + "/" + thn
	var pts = {
		"nama": $("span#rspnamapts").html(),
		"umur": $("span#rspumur").html(),
		"berat": $("span#rspbb").html(),
		"alamat": $("span#rspalmt").html(),
		"alergi": $("span#rspalergi").html(),
		"diag": $("span#rspdiag").html()
	}
	// var listini = $(this).children()
	var obat = []
	var puyer = []
	$(".listobat").each(function(){
		// console.log($(this).find('.obatpuyer'))
		if ($(this).find('.obatpuyer').length == 0){
			// console.log("not puyer")
			// var namaobat = $(this).find(".isianobat").val();
			// var jumlah = $(this).find('.jum-obat').val();
			// var instruksi = $(this).find('.instruksi').val();
			// var keterangan = $(this).find('.keterangan').val();
			arr = {
				"obat": $(this).find(".isianobat").val(),
				"jumlah": $(this).find('.jum-obat').val(),
				"instruksi": $(this).find('.instruksi').val(),
				"keterangan": $(this).find('.keterangan').val()
				}
			obat.push(arr)

		}else{
			// console.log("puyer")
			pyr = []
			$(this).find('.list-puyer').not('.template-puyer').each(function(){
				arr = {
					"obat": $(this).find('.obat-puyer').val(),
					"takaran": $(this).find('.takaran-obat').val()
				}
				pyr.push(arr)
			})
			mix = {
				"satuobat": pyr,
				"racikan": $(this).find('.peracikan').val(),
				"jml-racikan": $(this).find('.jml-racikan').val(),
				"instruksi": $(this).find('.instruksi').val(),
				"keterangan": $(this).find('.keterangan').val()
			}
			puyer.push(mix)
		}
	})
		send = {
			"dokter": $("#dokter").html(),
			"tanggal": strDate,
			"listobat": obat,
			"listpuyer": puyer,
			"pasien": pts
		}
		console.log(JSON.stringify(send))
		$.post("buatresep", 
		{"send" : JSON.stringify(send)},
		function(data){
			let pdfWindow = window.open("")
			pdfWindow.document.write("<iframe width='100%' height='100%' src='data:application/pdf;base64, " + encodeURI(data)+"'></iframe>")
		})

	});

});


