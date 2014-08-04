part of components.{{.Name}}s;
@Component(selector: 'edit',  templateUrl: 'packages/{{.App.Name}}/components/{{.Name}}s/form.html', publishAs: 'C')
class Edit extends Cntr{
  {{.Cname}} record;
  bool formSubmitted;

  Edit(RouteProvider routeProvider,Db db,NgRoutingHelper locationService) :
    super({{.Cname}},routeProvider,db,locationService) {
    ready.then((_){
      //store.once();
      record = store.find(int.parse(Id, onError: (_) => 0));
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
    locationService.router.gotoUrl('/');
  }
}
