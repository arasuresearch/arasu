part of {{.App.Name}}.controllers.{{.Name}}s;
@Component(selector: '{{.Name}}s-new',  
  templateUrl: 'packages/{{.App.Name}}/views/{{.Name}}s/form.html')
class {{.Cname}}sNew extends Cntr {
  {{.Cname}} record;
  bool formSubmitted;
  NgForm form;
  
  {{.Cname}}sNew(RouteProvider routeProvider,Db db,NgRoutingHelper locationService) :
    super({{.Cname}},routeProvider,db,locationService) {
    ready.then((_){
      record = store.New();
      record.updateClone();
    }).catchError((e)=>Logger.root.shout(e));
  }

  void submit(form) {
    formSubmitted = true;
    if(form.invalid) return;
    if (record.hasChanged) {
      store.add(record).then((_){
        locationService.router.gotoUrl('/{{.Name}}s');
      }).catchError((e)=>Logger.root.shout(e));
    }
  }
}