import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';
import { CheckInCountList } from '../../../api/models/check-in';
import { ApiService } from '../../../api/api.service';

@Component({
  selector: 'anms-check-in-count',
  templateUrl: './check-in-count.component.html',
  styleUrls: ['./check-in-count.component.scss'],
  changeDetection: ChangeDetectionStrategy.Default
})
export class CheckInCountComponent implements OnInit {
  public checkInCountList: CheckInCountList | undefined;

  constructor(
    private api: ApiService,
  ) { }

  ngOnInit(): void {
    this.api.getCheckInCountTop().subscribe(data => this.checkInCountList = data)
  }
}
