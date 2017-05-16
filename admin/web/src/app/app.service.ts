import {Injectable} from '@angular/core';
import {StoreService} from 'le5le-store';

import { HttpService } from '../core/http.service';

@Injectable()
export class AppService {
  public constructor(protected http: HttpService, protected store: StoreService) {
  }

  GetSignConfig(): Promise<any> {
    return this.http.Get('/api/sign/config');
  }

}
