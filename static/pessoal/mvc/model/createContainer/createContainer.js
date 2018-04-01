golbalOnInit.push(function(){
  $("#createAContainerIconTextButton").kendoButton({
    icon: "plus",
    click: function(e){
      $("#windowCreateContainer").data("kendoWindow").open().center();
    }
  });
  
  $("#windowCreateContainer").kendoWindow({
    title: translator.get("Create a container"),
    //data: {id: ""},
    modal: true,
    visible: false,
    resizable: true,
    width: 700,
    height: 400,
    actions: [
      "Maximize",
      "Close"
    ],
    close: function(){

    },
    open: function(){

    }
  }).data("kendoWindow").center();
});