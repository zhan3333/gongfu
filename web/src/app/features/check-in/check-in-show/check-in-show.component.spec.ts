import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CheckInShowComponent } from './check-in-show.component';

describe('CheckInShowComponent', () => {
  let component: CheckInShowComponent;
  let fixture: ComponentFixture<CheckInShowComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
    imports: [CheckInShowComponent]
})
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(CheckInShowComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
