import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';
import { CheckInCountList } from '../../../api/models/check-in';
import { ApiService } from '../../../api/api.service';

@Component({
  selector: 'anms-check-in-continuous',
  templateUrl: './check-in-continuous.component.html',
  styleUrls: ['./check-in-continuous.component.scss'],
  changeDetection: ChangeDetectionStrategy.Default
})
export class CheckInContinuousComponent implements OnInit {
  public checkInCountList: CheckInCountList | undefined;

  constructor(
    private api: ApiService,
  ) { }

  ngOnInit(): void {
    this.api.getCheckInContinuousTop().subscribe(data => this.checkInCountList = data)
  }
}
