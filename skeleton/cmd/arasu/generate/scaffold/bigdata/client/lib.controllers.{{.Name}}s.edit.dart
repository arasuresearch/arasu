part of {{.App.Name}}.controllers.{{.Name}}s;
@Component(selector: '{{.Name}}s-edit',  
  templateUrl: 'packages/{{.App.Name}}/views/{{.Name}}s/form.html')
class {{.Cname}}sEdit extends Cntr{
  {{.Cname}} record;
  bool formSubmitted;
  NgForm form;

  {{.Cname}}sEdit(RouteProvider routeProvider,Db db,NgRoutingHelper locationService) :
    super({{.Cname}},routeProvider,db,locationService) {
    ready.then((_){
      //store.once();
      record = store.find(Id);//store.find(int.parse(Id, onError: (_) => 0));
    }).catchError((e)=>Logger.root.shout(e));
  }

//  id = int.parse(routeProvider.parameters['id'], onError: (_) => 0);

  void submit(form) {
    formSubmitted = true;
    if(form.invalid) return;
    if (record.hasChanged)  {
      store.put(record).then((_){
      })
      ..catchError((e)=>Logger.root.shout(e));
    }
    locationService.router.gotoUrl('/{{.Name}}s');
  }
}
