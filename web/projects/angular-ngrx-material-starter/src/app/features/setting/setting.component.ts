import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';
import { User } from '../../api/models/user';
import { ApiService } from '../../api/api.service';

@Component({
  selector: 'anms-setting',
  templateUrl: './setting.component.html',
  styleUrls: ['./setting.component.scss'],
  changeDetection: ChangeDetectionStrategy.Default
})
export class SettingComponent implements OnInit {

  public user: User | undefined

  constructor(
    private api: ApiService
  ) { }

  ngOnInit(): void {
    this.api.me().subscribe(user => this.user = user)
  }
}
