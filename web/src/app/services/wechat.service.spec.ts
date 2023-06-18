import { TestBed } from '@angular/core/testing';

import { WechatService } from './wechat.service';

describe('WechatService', () => {
  let service: WechatService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(WechatService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
