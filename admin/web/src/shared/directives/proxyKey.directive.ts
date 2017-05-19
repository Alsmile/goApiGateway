import { Directive, Input } from '@angular/core';
import { NG_VALIDATORS, Validator, AbstractControl } from '@angular/forms';

@Directive({
  selector: '[proxyKey]',
  providers: [{
    provide: NG_VALIDATORS,
    useExisting: ProxyKeyValidator,
    multi: true
  }]
})
export class ProxyKeyValidator implements Validator {
  validate(c: AbstractControl): {[key: string]: any} {
    if (!c.value) return;

    if (c.value[0] !== '/' || c.value.length < 2 || c.value.indexOf('/', 1) > -1) {
      return {'proxyKey': true};
    }
  }
}
